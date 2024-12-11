-- MariaDB dump 10.19-11.2.2-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: ops2
-- ------------------------------------------------------
-- Server version	11.2.2-MariaDB-1:11.2.2+maria~ubu2204

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES
(1,'2024-07-09 10:58:54.000','2024-12-10 17:02:22.252',NULL,1,'admin','$2a$12$G3uRulMbELE9geIbsB0tDOQP.WraVHD94M3uYWwGGgZr4uZ77EIZW','1','管理员','','1234567890@gmail.com'),
(3,'2024-07-10 14:02:54.000','2024-09-27 14:41:19.628','2024-12-11 16:07:53.000',1,'normal1','$2a$12$KbD7XBnbI3uZF60hpl3VrehcrxUl0MKLhHS5lhzzznUHYdsH6zo9i','2','普通用户1','','123456789@qq.com');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `user_role`
--

LOCK TABLES `user_role` WRITE;
/*!40000 ALTER TABLE `user_role` DISABLE KEYS */;
INSERT INTO `user_role` VALUES
(1,1);
/*!40000 ALTER TABLE `user_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `role`
--

LOCK TABLES `role` WRITE;
/*!40000 ALTER TABLE `role` DISABLE KEYS */;
INSERT INTO `role` VALUES
(1,'2024-07-09 11:01:06.000','2024-12-10 17:02:22.000',NULL,'管理员','ADMIN','管理员权限(最高权限)'),
(2,'2024-07-10 15:24:53.000','2024-09-27 14:41:19.000','2024-12-11 16:07:57.000','普通用户组','NORMAL','普通用户权限');
/*!40000 ALTER TABLE `role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `menu`
--

LOCK TABLES `menu` WRITE;
/*!40000 ALTER TABLE `menu` DISABLE KEYS */;
INSERT INTO `menu` VALUES
(1,'0000-00-00 00:00:00.000','2024-07-15 15:03:36.000',NULL,1,0,2,'首页','home','/home','layout.base$view.home',1,'route.home','mdi:monitor-dashboard',1,0,0,NULL,'0',0,0,NULL,NULL,NULL),
(2,'2024-07-11 14:45:59.000','2024-09-29 16:52:24.000',NULL,1,0,1,'功能','function','/function','layout.base',2,'route.function','icon-park-outline:all-application',1,0,1,NULL,'0',0,0,NULL,NULL,NULL),
(3,'0000-00-00 00:00:00.000','2024-07-15 15:49:41.000',NULL,1,2,2,'多标签页','function_multi-tab','/function/multi-tab','view.function_multi-tab',1,'route.function_multi-tab','ic:round-tab',1,1,1,0,'0',0,0,NULL,NULL,NULL),
(4,'2024-07-11 16:05:43.000','2024-07-11 16:05:43.000',NULL,1,2,2,'标签页','function_tab','/function/tab','view.function_tab',2,'route.function_tab','ic:round-tab',1,0,1,NULL,'0',0,0,NULL,NULL,NULL),
(5,'2024-07-11 16:49:02.000','2024-07-11 16:49:02.000',NULL,1,0,1,'异常页','exception','/exception','layout.base',3,'route.exception','ant-design:exception-outlined',1,0,1,NULL,'0',1,0,NULL,NULL,NULL),
(6,'2024-07-11 17:02:36.000','2024-07-11 17:02:36.000',NULL,1,5,2,'403','exception_403','/exception/403','view.403',1,'route.exception_403','ic:baseline-block',1,0,1,NULL,'0',1,0,NULL,NULL,NULL),
(7,'2024-07-11 17:13:35.000','2024-07-11 17:13:35.000',NULL,1,5,2,'404','exception_404','/exception/404','view.404',2,'route.exception_404','ic:baseline-web-asset-off',1,0,1,NULL,'0',1,0,NULL,NULL,NULL),
(8,'2024-07-11 17:16:37.000','2024-07-11 17:16:37.000',NULL,1,5,2,'500','exception_500','/exception/500','view.500',3,'route.exception_500','ic:baseline-wifi-off',1,0,1,NULL,'0',1,0,NULL,NULL,NULL),
(9,'2024-07-11 17:50:50.000','2024-10-21 15:48:08.000',NULL,1,0,1,'系统管理','manage','/manage','layout.base',8,'route.manage','carbon:cloud-service-management',1,0,0,NULL,'0',0,0,NULL,NULL,NULL),
(10,'2024-07-11 19:33:45.000','2024-07-15 15:49:41.000',NULL,1,9,2,'用户管理','manage_user','/manage/user','view.manage_user',1,'route.manage_user','ic:round-manage-accounts',1,0,0,NULL,'0',0,0,NULL,NULL,NULL),
(11,'2024-07-11 19:38:52.000','2024-07-15 15:49:41.000',NULL,1,9,2,'角色管理','manage_role','/manage/role','view.manage_role',2,'route.manage_role','carbon:user-role',1,0,0,NULL,'0',0,0,NULL,NULL,NULL),
(12,'2024-07-11 19:40:01.000','2024-08-15 14:14:54.000',NULL,1,9,2,'菜单管理','manage_menu','/manage/menu','view.manage_menu',3,'route.manage_menu','material-symbols:route',1,0,0,NULL,'1',0,0,NULL,NULL,NULL),
(13,'2024-07-11 20:48:28.000','2024-08-30 17:58:25.000',NULL,1,9,2,'用户详情','manage_user-detail','/manage/user-detail/:id','view.manage_user-detail',4,'route.manage_user-detail','bxs:user-detail',1,0,1,0,'0',0,0,'true',NULL,NULL),
(16,'2024-08-15 20:19:34.000','2024-08-15 20:19:34.000',NULL,1,0,2,'内嵌页面','iframe-page','/iframe-page/:url','layout.base$view.iframe-page',5,'route.iframe-page',NULL,1,0,1,NULL,'1',1,0,'true',NULL,NULL),
(17,'2024-08-16 16:46:58.000','2024-08-16 16:46:58.000',NULL,1,0,2,'个人中心','user-center','/user-center','layout.base$view.user-center',6,'route.user-center',NULL,1,0,1,NULL,'0',0,0,NULL,NULL,NULL),
(18,'2024-08-19 11:30:38.000','2024-08-19 11:30:38.000',NULL,1,0,2,'登录','login','/login/:module(pwd-login|code-login|register|reset-pwd|bind-wechat)?','layout.blank$view.login',7,'route.login',NULL,1,0,1,NULL,'0',1,0,'true',NULL,NULL),
(19,'2024-09-05 15:00:00.000','2024-09-05 15:28:26.000',NULL,1,9,2,'API管理','manage_api','/manage/api','view.manage_api',5,'route.manage_api','hugeicons:api',1,0,0,NULL,'0',0,0,NULL,NULL,NULL),
(20,'2024-09-06 15:57:19.000','2024-09-06 15:57:32.000',NULL,1,9,2,'用户操作记录','manage_user-record','/manage/user-record','view.manage_user-record',6,'route.manage_user-record','octicon:log-24',1,0,0,NULL,'0',0,0,NULL,NULL,NULL),
(21,'2024-09-24 11:03:40.000','2024-10-21 15:48:13.000',NULL,1,0,1,'资产管理','asset','/asset','layout.base',9,'route.asset','fluent-mdl2:fixed-asset-management',1,0,0,NULL,'0',0,0,NULL,NULL,NULL),
(22,'2024-09-24 11:09:44.000','2024-09-24 11:10:17.000',NULL,1,21,2,'项目管理','asset_project','/asset/project','view.asset_project',1,'route.asset_project','arcticons:projectm',1,0,0,NULL,'0',0,0,NULL,NULL,NULL),
(23,'2024-09-27 14:47:51.000','2024-09-27 14:47:51.000',NULL,1,21,2,'服务器管理','asset_host','/asset/host','view.asset_host',2,'route.asset_host','clarity:host-group-solid',1,0,0,NULL,'0',0,0,NULL,NULL,NULL),
(24,'2024-09-29 14:01:28.000','2024-09-29 14:08:18.000',NULL,1,0,1,'文档','document','/document','layout.base',4,'route.document','et:document',1,0,0,NULL,'0',0,0,NULL,NULL,NULL),
(25,'2024-09-29 17:26:20.000','2024-09-29 17:52:45.000',NULL,1,24,2,'项目文档(内链)','document_project','/document/project','view.iframe-page',1,'route.document_project','gala:file-document',1,0,0,NULL,'0',0,0,'{\"url\":\"https://github.com/1113464192/FreeOpsServer\"}',NULL,NULL),
(26,'2024-09-29 17:51:00.000','2024-09-29 17:51:00.000',NULL,1,24,2,'项目文档(外链)','document_project-link','/document/project-link','',2,'route.document_project-link','gala:file-document',1,0,0,NULL,'0',0,0,NULL,NULL,'https://github.com/1113464192/FreeOpsClient'),
(27,'2024-10-21 15:50:29.000','2024-10-21 15:50:29.000',NULL,1,0,1,'运维操作管理','ops-manage','/ops-manage','layout.base',10,'route.ops-manage','carbon:operations-field',1,0,0,NULL,'0',0,0,NULL,NULL,NULL),
(28,'2024-10-21 15:57:56.000','2024-10-21 15:57:56.000',NULL,1,27,2,'游戏管理','ops-manage_game','/ops-manage/game','view.ops-manage_game',1,'route.ops-manage_game','tabler:businessplan',1,0,0,NULL,'0',0,0,NULL,NULL,NULL),
(29,'2024-11-21 13:38:31.000','2024-11-21 13:38:31.000',NULL,1,27,2,'操作模板管理','ops-manage_template','/ops-manage/template','view.ops-manage_template',2,'route.ops-manage_template','tdesign:template',1,0,0,NULL,'0',0,0,NULL,NULL,NULL),
(30,'2024-11-26 16:57:56.000','2024-11-26 16:57:56.000',NULL,1,27,2,'参数模板管理','ops-manage_param-template','/ops-manage/param-template','view.ops-manage_param-template',3,'route.ops-manage_param-template','arcticons:param',1,0,0,NULL,'0',0,0,NULL,NULL,NULL),
(31,'2024-11-29 10:35:03.000','2024-11-29 10:35:03.000',NULL,1,27,2,'任务管理','ops-manage_task','/ops-manage/task','view.ops-manage_task',4,'route.ops-manage_task','tdesign:task',1,0,0,NULL,'0',0,0,NULL,NULL,NULL),
(32,'2024-12-03 14:56:02.000','2024-12-03 14:56:53.000',NULL,1,27,2,'审批任务','ops-manage_approve-task','/ops-manage/approve-task','view.ops-manage_approve-task',5,'route.ops-manage_approve-task','material-symbols-light:order-approve-outline-rounded',1,0,0,NULL,'0',0,0,NULL,NULL,NULL),
(33,'2024-12-05 11:23:56.000','2024-12-05 11:24:31.000',NULL,1,27,2,'任务日志','ops-manage_task-log','/ops-manage/task-log','view.ops-manage_task-log',6,'route.ops-manage_task-log','mdi:math-log',1,0,0,NULL,'0',0,0,NULL,NULL,NULL);
/*!40000 ALTER TABLE `menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `menu_role`
--

LOCK TABLES `menu_role` WRITE;
/*!40000 ALTER TABLE `menu_role` DISABLE KEYS */;
INSERT INTO `menu_role` VALUES
(1,1),
(2,1),
(3,1),
(4,1),
(5,1),
(6,1),
(7,1),
(8,1),
(9,1),
(10,1),
(11,1),
(12,1),
(13,1),
(16,1),
(17,1),
(18,1),
(19,1),
(20,1),
(21,1),
(22,1),
(23,1),
(24,1),
(25,1),
(26,1),
(27,1),
(28,1),
(29,1),
(30,1),
(31,1),
(32,1),
(33,1);
/*!40000 ALTER TABLE `menu_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `api`
--

LOCK TABLES `api` WRITE;
/*!40000 ALTER TABLE `api` DISABLE KEYS */;
INSERT INTO `api` VALUES
(1,'2024-07-10 14:22:47.000','2024-09-03 11:12:17.000',NULL,'/users/ssh-key','PUT','users','用户私钥'),
(2,'2024-07-10 14:23:06.000','2024-09-09 14:40:40.000',NULL,'/users/password','PATCH','users','修改用户密码'),
(3,'2024-07-10 14:23:22.000','2024-07-10 16:27:24.000',NULL,'/users','DELETE','users','删除用户'),
(4,'2024-07-10 14:23:34.000','2024-09-09 14:40:23.000',NULL,'/users/privilege','GET','users','查询用户权限'),
(5,'2024-07-10 14:28:19.000','2024-09-09 15:05:48.000',NULL,'/users/logout','POST','users','登出'),
(8,'2024-08-20 14:23:00.000','2024-08-20 14:22:54.000',NULL,'/menus/user-routes','GET','menus','获取用户路由'),
(9,'2024-09-05 20:32:25.000','2024-09-05 20:32:25.000',NULL,'/auth/login','POST','auth','用户登录'),
(10,'2024-09-05 20:38:52.000','2024-09-05 20:38:52.000',NULL,'/auth/refreshToken','POST','auth','热刷Token'),
(11,'2024-09-09 14:38:30.000','2024-09-09 14:39:07.000',NULL,'/users','POST','users','新增/修改用户'),
(12,'2024-09-09 14:39:36.000','2024-09-09 14:39:36.000',NULL,'/users','GET','users','查询用户切片'),
(13,'2024-09-09 15:06:19.000','2024-09-09 15:06:19.000',NULL,'/users/history-action','GET','users','查询用户所有的历史操作'),
(14,'2024-09-09 15:06:41.000','2024-09-09 15:06:41.000',NULL,'/users/history-month-exist','GET','users','查询有多少个月份表可供查询'),
(15,'2024-09-09 15:26:32.000','2024-09-09 15:26:32.000',NULL,'/users/bind-roles','PUT','users','用户绑定角色'),
(16,'2024-09-09 15:27:12.000','2024-09-09 15:27:12.000',NULL,'/users/roles','GET','users','查看用户所有角色'),
(17,'2024-09-09 15:33:27.000','2024-09-09 15:33:27.000',NULL,'/roles','POST','roles','新增/修改角色'),
(18,'2024-09-09 15:35:25.000','2024-09-09 15:35:25.000',NULL,'/roles','GET','roles','获取角色列表'),
(19,'2024-09-09 15:37:30.000','2024-09-09 15:37:30.000',NULL,'/roles/all-summary','GET','roles','获取所有角色的简略信息'),
(20,'2024-09-09 15:37:50.000','2024-09-09 15:37:50.000',NULL,'/roles','DELETE','roles','删除角色'),
(21,'2024-09-09 15:40:23.000','2024-09-09 15:40:23.000',NULL,'/roles/bind','PUT','roles','角色绑定关系'),
(22,'2024-09-09 15:42:10.000','2024-09-09 15:42:10.000',NULL,'/roles/menus','GET','roles','获取角色的菜单'),
(23,'2024-09-09 15:44:00.000','2024-09-09 15:44:00.000',NULL,'/roles/apis','GET','roles','获取角色的API'),
(24,'2024-09-09 15:47:18.000','2024-09-09 15:47:18.000',NULL,'/roles/buttons','GET','roles','获取角色的按钮'),
(25,'2024-09-09 15:47:37.000','2024-09-09 15:47:37.000',NULL,'/roles/users','GET','roles','获取角色绑定的用户'),
(26,'2024-09-09 15:49:44.000','2024-09-09 15:49:44.000',NULL,'/buttons','POST','buttons','新增/修改按钮'),
(27,'2024-09-09 15:49:55.000','2024-09-09 15:49:55.000',NULL,'/buttons','GET','buttons','获取按钮列表'),
(28,'2024-09-09 15:50:14.000','2024-12-10 17:47:24.000',NULL,'/buttons/menus','DELETE','buttons','删除菜单的按钮,提供菜单ID'),
(29,'2024-09-09 15:51:53.000','2024-09-09 15:51:53.000',NULL,'/menus','POST','menus','新增/修改菜单'),
(30,'2024-09-09 15:52:13.000','2024-09-09 15:52:13.000',NULL,'/menus','GET','menus','获取菜单信息'),
(31,'2024-09-09 15:56:26.000','2024-09-09 15:56:26.000',NULL,'/menus','DELETE','menus','删除菜单'),
(32,'2024-09-09 15:56:47.000','2024-09-09 15:56:47.000',NULL,'/menus/buttons','GET','menus','获取菜单下所有按钮'),
(33,'2024-09-09 15:57:08.000','2024-09-09 15:57:08.000',NULL,'/menus/all-pages','GET','menus','获取所有页面'),
(34,'2024-09-09 15:57:29.000','2024-09-09 15:57:29.000',NULL,'/menus/tree','GET','menus','获取菜单树'),
(35,'2024-09-09 15:57:59.000','2024-09-09 15:57:59.000',NULL,'/apis','POST','apis','新增/修改API'),
(36,'2024-09-09 15:58:14.000','2024-09-09 15:58:14.000',NULL,'/apis','GET','apis','获取API列表'),
(37,'2024-09-09 15:58:28.000','2024-09-09 15:58:28.000',NULL,'/apis','DELETE','apis','删除API'),
(38,'2024-09-09 15:58:49.000','2024-09-09 15:58:49.000',NULL,'/apis/group','GET','apis','获取存在的API组'),
(39,'2024-09-09 16:00:49.000','2024-09-09 16:00:49.000',NULL,'/apis/tree','GET','apis','获取API树'),
(40,'2024-11-18 11:32:51.000','2024-11-18 11:32:51.000',NULL,'/ops/approve-task','PUT','ops','运维操作相关'),
(41,'2024-12-10 17:10:41.000','2024-12-10 17:12:13.000',NULL,'/auth/error','GET','auth','自定义错误接口'),
(42,'2024-12-10 17:11:13.000','2024-12-10 17:12:24.000',NULL,'/auth/constant-routes','GET','auth','获取所有常量路由'),
(43,'2024-12-10 17:11:59.000','2024-12-10 17:11:59.000',NULL,'/ops/task-need-approve','GET','ops','查询用户是否有待审批的任务'),
(44,'2024-12-10 17:12:49.000','2024-12-10 17:12:49.000',NULL,'/ops/task-running-ws','GET','ops','实时查看运行中的任务'),
(45,'2024-12-10 17:18:28.000','2024-12-10 17:18:28.000',NULL,'/users/project-options','GET','users','查看用户所有项目选项'),
(46,'2024-12-10 17:41:25.000','2024-12-10 17:41:25.000',NULL,'/roles/projects','GET','roles','获取角色绑定的项目'),
(47,'2024-12-10 17:54:37.000','2024-12-10 17:54:37.000',NULL,'/menus/is-route-exist','GET','menus','判断路由是否存在'),
(48,'2024-12-10 17:56:29.000','2024-12-10 17:56:29.000',NULL,'/projects','POST','projects','新增/修改项目'),
(49,'2024-12-10 17:57:08.000','2024-12-10 17:57:08.000',NULL,'/projects','GET','projects','查询项目'),
(50,'2024-12-10 19:53:57.000','2024-12-10 19:53:57.000',NULL,'/projects/all-summary','GET','projects','获取项目列表'),
(51,'2024-12-10 19:54:20.000','2024-12-10 19:54:20.000',NULL,'/projects','DELETE','projects','删除项目'),
(52,'2024-12-10 19:54:39.000','2024-12-10 19:54:39.000',NULL,'/projects/hosts','GET','projects','查询项目关联的服务器'),
(53,'2024-12-10 19:54:52.000','2024-12-10 19:54:52.000',NULL,'/projects/games','GET','projects','查询项目关联的游戏'),
(54,'2024-12-10 19:55:06.000','2024-12-10 19:55:06.000',NULL,'/projects/assets-total','GET','projects','查询项目各资产总数'),
(55,'2024-12-10 19:55:43.000','2024-12-10 19:55:43.000',NULL,'/hosts','POST','hosts','新增/修改服务器'),
(56,'2024-12-10 19:57:13.000','2024-12-10 19:57:13.000',NULL,'/hosts','GET','hosts','查询服务器'),
(57,'2024-12-10 19:57:29.000','2024-12-10 19:57:29.000',NULL,'/hosts','DELETE','hosts','删除服务器'),
(58,'2024-12-10 19:57:41.000','2024-12-10 19:57:41.000',NULL,'/hosts/summary','GET','hosts','获取服务器列表'),
(59,'2024-12-10 19:57:54.000','2024-12-10 19:57:54.000',NULL,'/hosts/game-info','GET','hosts','获取服务器的游戏信息'),
(60,'2024-12-10 19:58:22.000','2024-12-10 19:58:22.000',NULL,'/ops/template','POST','ops','创建/修改 模板'),
(61,'2024-12-10 20:00:06.000','2024-12-10 20:00:06.000',NULL,'/ops/template','GET','ops','查看模板'),
(62,'2024-12-10 20:00:20.000','2024-12-10 20:00:20.000',NULL,'/ops/template','DELETE','ops','删除模板'),
(63,'2024-12-10 20:04:56.000','2024-12-10 20:04:56.000',NULL,'/ops/param-template','POST','ops','创建/修改 获取参数模板 (从运营文案信息获取参数的正则模板)'),
(64,'2024-12-10 20:05:16.000','2024-12-10 20:05:16.000',NULL,'/ops/param-template','GET','ops','查看参数模板'),
(65,'2024-12-10 20:05:28.000','2024-12-10 20:05:28.000',NULL,'/ops/param-template','DELETE','ops','删除参数模板'),
(66,'2024-12-10 20:05:57.000','2024-12-10 20:05:57.000',NULL,'/ops/bind-template-params','PUT','ops','绑定模板参数'),
(67,'2024-12-11 11:07:53.000','2024-12-11 11:07:53.000',NULL,'/ops/template-params','GET','ops','查看模板关联的参数'),
(68,'2024-12-11 11:08:13.000','2024-12-11 11:08:13.000',NULL,'/ops/task','POST','ops','创建/修改 任务(拼接执行模板顺序的任务)'),
(69,'2024-12-11 11:08:39.000','2024-12-11 11:08:39.000',NULL,'/ops/task','GET','ops','查看任务'),
(70,'2024-12-11 11:08:56.000','2024-12-11 11:08:56.000',NULL,'/ops/task','DELETE','ops','删除任务'),
(71,'2024-12-11 11:09:23.000','2024-12-11 11:09:23.000',NULL,'/ops/commands','POST','ops','查看根据参数会生成的命令'),
(72,'2024-12-11 11:09:58.000','2024-12-11 11:09:58.000',NULL,'/ops/submit-task','POST','ops','提交任务'),
(73,'2024-12-11 11:11:15.000','2024-12-11 11:11:15.000',NULL,'/ops/task-pending','GET','ops','查询用户待审批的任务'),
(74,'2024-12-11 11:11:32.000','2024-12-11 11:11:32.000',NULL,'/ops/run-task-check-script','POST','ops','执行并等待运营检查脚本返回结果'),
(75,'2024-12-11 11:11:50.000','2024-12-11 11:11:50.000',NULL,'/ops/task-log','GET','ops','查看任务日志'),
(76,'2024-12-11 11:13:10.000','2024-12-11 11:13:10.000',NULL,'/clouds/create/project','POST','clouds','创建云项目'),
(77,'2024-12-11 11:13:24.000','2024-12-11 11:13:24.000',NULL,'/clouds/create/host','POST','clouds','创建云服务器'),
(78,'2024-12-11 11:13:50.000','2024-12-11 11:13:50.000',NULL,'/clouds/update/project','POST','clouds','更新云项目'),
(79,'2024-12-11 11:14:03.000','2024-12-11 11:14:03.000',NULL,'/clouds/query/project','GET','clouds','查询云项目ID');
/*!40000 ALTER TABLE `api` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-12-11  8:16:11
