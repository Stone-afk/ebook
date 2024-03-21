package service

import (
	"context"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/user/domain"
	"ebook/cmd/user/events"
	"ebook/cmd/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	producer events.SyncSearchEventProducer
	repo     repository.UserRepository
	l        logger.Logger
}

func NewUserService(producer events.SyncSearchEventProducer, repo repository.UserRepository, l logger.Logger) UserService {
	return &userService{
		producer: producer,
		repo:     repo,
		l:        l,
	}
}

func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	// 写法1
	// 这种是简单的写法，依赖与 Web 层保证没有敏感数据被修改
	// 也就是说，你的基本假设是前端传过来的数据就是不会修改 Email，Phone 之类的信息的。
	//return svc.repo.Update(ctx, user)

	// 写法2
	// 这种是复杂写法，依赖于 repository 中更新会忽略 0 值
	// 这个转换的意义在于，你在 service 层面上维护住了什么是敏感字段这个语义
	user.Email = ""
	user.Phone = ""
	user.Password = ""
	return svc.repo.Update(ctx, user)
}

func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

// FindOrCreate 如果手机号不存在，那么会初始化一个用户
func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	// 这是一种优化写法, 大部分人会命中这个分支
	if err != repository.ErrUserNotFound {
		return u, err
	}
	// 这里，把 phone 脱敏之后打出来
	//zap.L().Info("用户未注册", zap.String("phone", phone))
	//svc.logger.Info("用户未注册", zap.String("phone", phone))
	svc.l.Info("用户未注册", logger.String("phone", phone))
	// 要执行注册
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	// 注册有问题，但是又不是用户手机号码冲突，说明是系统错误
	if err != nil && err != repository.ErrUserDuplicate {
		return domain.User{}, err
	}
	u, err = svc.repo.FindByPhone(ctx, phone)
	if err == nil {
		go func() {
			svc.sendSyncEventToSearch(ctx, u)
		}()
	}
	// 主从模式下，这里要从主库中读取，暂时我们不需要考虑
	return u, err
}

func (svc *userService) sendSyncEventToSearch(ctx context.Context, u domain.User) {
	evt := events.UserEvent{
		Id:       u.Id,
		Email:    u.Email,
		Phone:    u.Phone,
		Nickname: u.Nickname,
	}
	er := svc.producer.ProduceSyncEvent(ctx, evt)
	if er != nil {
		svc.l.Error("ProduceSyncEvent 发送同步搜索用户事件失败", logger.Error(er))
		er = svc.producer.ProduceStandardSyncEvent(ctx, evt)
		if er != nil {
			svc.l.Error("ProduceStandardSyncEvent 发送同步搜索用户事件失败", logger.Error(er))
		}
	}
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, wechatInfo.OpenId)
	// 这是一种优化写法, 大部分人会命中这个分支
	if err != repository.ErrUserNotFound {
		return u, err
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	// 因为这里会遇到主从延迟的问题
	u, err = svc.repo.FindByWechat(ctx, wechatInfo.OpenId)
	if err == nil {
		go func() {
			svc.sendSyncEventToSearch(ctx, u)
		}()
	}
	return u, err
}

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, err
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}
