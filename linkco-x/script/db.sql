
--切换数据库
use xworld;
 
-- Table structure for table `user`
DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `id` varchar(50) NOT NULL,
  `username` varchar(50) NOT NULL DEFAULT "xxx",
  `basic_data` BLOB NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
