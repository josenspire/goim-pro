-- 创建数据库:goim
-- create database goim;

use goim;

-- 创建用户: goimroot
-- CREATE USER 'goimroot'@'%' IDENTIFIED BY 'Password1!';
-- 授予用户通过外网IP对数据库 “goim” 的全部权限
-- all 可以替换为 select,delete,update,create,drop
grant all privileges on *.* to 'root'@'%';

grant all privileges on *.* to 'goimroot'@'%';

-- 刷新权限
flush privileges;
