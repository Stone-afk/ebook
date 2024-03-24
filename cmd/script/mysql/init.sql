create database ebook;
create database ebook_interactive;
create database ebook_article;
create database ebook_user;
create database ebook_payment;
create database ebook_account;
create database ebook_reward;
create database ebook_comment;
create database ebook_tag;

# 准备 canal 用户
CREATE USER 'canal'@'%' IDENTIFIED BY 'canal';
GRANT ALL PRIVILEGES ON *.* TO 'canal'@'%' WITH GRANT OPTION;