package intergration

import (
	"context"
	followv1 "ebook/cmd/api/proto/gen/followrelation/v1"
	"ebook/cmd/followrelation/intergration/startup"
	"ebook/cmd/followrelation/repository/dao"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type FollowRelationSuite struct {
	suite.Suite
	db     *gorm.DB
	rdb    redis.Cmdable
	server followv1.FollowServiceServer
}

func (s *FollowRelationSuite) SetupSuite() {
	s.db = startup.InitTestDB()
	s.rdb = startup.InitRedis()
	s.server = startup.InitServer()
}

func (s *FollowRelationSuite) TearDownSuite() {
	err := s.db.Where("id > ?", 0).Delete(&dao.FollowRelation{}).Error
	require.NoError(s.T(), err)
}

func (s *FollowRelationSuite) TestFollowRelation_ADD() {
	testcases := []struct {
		name    string
		before  func()
		req     *followv1.FollowRequest
		wantVal *followv1.FollowRelation
		wantErr error
	}{
		{
			name: "添加正常",
			before: func() {
			},
			req: &followv1.FollowRequest{
				Followee: 1,
				Follower: 2,
			},
			wantVal: &followv1.FollowRelation{
				Followee: 1,
				Follower: 2,
			},
		},
		{
			name: "关注关系重复",
			before: func() {
				_, err := s.server.Follow(context.Background(), &followv1.FollowRequest{
					Followee: 2,
					Follower: 1,
				})
				require.NoError(s.T(), err)
			},
			req: &followv1.FollowRequest{
				Followee: 2,
				Follower: 1,
			},
			wantVal: &followv1.FollowRelation{
				Followee: 2,
				Follower: 1,
			},
		},
	}
	for _, tc := range testcases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before()
			_, err := s.server.Follow(context.Background(), tc.req)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			relation, err := s.GetFollowRelation(tc.req.Followee, tc.req.Follower)
			require.NoError(t, err)
			relation.Id = 0
			assert.Equal(t, tc.wantVal, relation)
		})
	}
}

func (s *FollowRelationSuite) TestFollowRelation_Info() {
	// 准备数据
	t := s.T()
	_, err := s.server.Follow(context.Background(), &followv1.FollowRequest{
		Followee: 8,
		Follower: 9,
	})
	require.NoError(t, err)
	relation, err := s.GetFollowRelation(8, 9)
	require.NoError(s.T(), err)
	resp, err := s.server.FollowInfo(context.Background(), &followv1.FollowInfoRequest{
		Follower: relation.Follower,
		Followee: relation.Followee,
	})
	require.NoError(t, err)
	assert.Equal(t, &followv1.FollowRelation{
		Id:       relation.Id,
		Followee: 8,
		Follower: 9,
	}, resp.FollowRelation)

}

func (s *FollowRelationSuite) GetFollowRelation(followee, follower int64) (*followv1.FollowRelation, error) {
	resp, err := s.server.FollowInfo(context.Background(), &followv1.FollowInfoRequest{
		Follower: follower,
		Followee: followee,
	})
	if err != nil {
		return nil, err
	}
	return resp.FollowRelation, nil
}

func TestFollowSuite(t *testing.T) {
	suite.Run(t, new(FollowRelationSuite))
}
