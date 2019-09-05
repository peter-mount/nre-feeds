package darwind3

import (
	"bytes"
	"encoding/xml"
	"github.com/peter-mount/nre-feeds/util"
	"log"
	"testing"
	"time"
)

const (
	issue6Xml      = "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><Pport xmlns=\"http://www.thalesgroup.com/rtti/PushPort/v16\" xmlns:ns2=\"http://www.thalesgroup.com/rtti/PushPort/Schedules/v3\" xmlns:ns3=\"http://www.thalesgroup.com/rtti/PushPort/Schedules/v2\" xmlns:ns4=\"http://www.thalesgroup.com/rtti/PushPort/Formations/v2\" xmlns:ns5=\"http://www.thalesgroup.com/rtti/PushPort/Forecasts/v3\" xmlns:ns6=\"http://www.thalesgroup.com/rtti/PushPort/Formations/v1\" xmlns:ns7=\"http://www.thalesgroup.com/rtti/PushPort/StationMessages/v1\" xmlns:ns8=\"http://www.thalesgroup.com/rtti/PushPort/TrainAlerts/v1\" xmlns:ns9=\"http://www.thalesgroup.com/rtti/PushPort/TrainOrder/v1\" xmlns:ns10=\"http://www.thalesgroup.com/rtti/PushPort/TDData/v1\" xmlns:ns11=\"http://www.thalesgroup.com/rtti/PushPort/Alarms/v1\" xmlns:ns12=\"http://thalesgroup.com/RTTI/PushPortStatus/root_1\" ts=\"2019-04-05T11:20:18.8416587+01:00\" version=\"16.0\"><uR updateOrigin=\"CIS\" requestSource=\"at01\" requestID=\"0000000000014065\"><trainOrder tiploc=\"BRGEND\" crs=\"BGN\" platform=\"1\"><ns9:set><ns9:first><ns9:trainID>6B09</ns9:trainID></ns9:first></ns9:set></trainOrder></uR></Pport>"
	trainOrderXML  = "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><Pport xmlns=\"http://www.thalesgroup.com/rtti/PushPort/v16\" xmlns:ns2=\"http://www.thalesgroup.com/rtti/PushPort/Schedules/v3\" xmlns:ns3=\"http://www.thalesgroup.com/rtti/PushPort/Schedules/v2\" xmlns:ns4=\"http://www.thalesgroup.com/rtti/PushPort/Formations/v2\" xmlns:ns5=\"http://www.thalesgroup.com/rtti/PushPort/Forecasts/v3\" xmlns:ns6=\"http://www.thalesgroup.com/rtti/PushPort/Formations/v1\" xmlns:ns7=\"http://www.thalesgroup.com/rtti/PushPort/StationMessages/v1\" xmlns:ns8=\"http://www.thalesgroup.com/rtti/PushPort/TrainAlerts/v1\" xmlns:ns9=\"http://www.thalesgroup.com/rtti/PushPort/TrainOrder/v1\" xmlns:ns10=\"http://www.thalesgroup.com/rtti/PushPort/TDData/v1\" xmlns:ns11=\"http://www.thalesgroup.com/rtti/PushPort/Alarms/v1\" xmlns:ns12=\"http://thalesgroup.com/RTTI/PushPortStatus/root_1\" ts=\"2019-04-05T15:09:59.3396587+01:00\" version=\"16.0\"><uR updateOrigin=\"CIS\" requestSource=\"AM07\" requestID=\"AM07573184\"><trainOrder tiploc=\"KENOLYM\" crs=\"KPA\" platform=\"2\"><ns9:set><ns9:first><ns9:rid wta=\"15:12\" wtd=\"15:12:30\" pta=\"15:11\" ptd=\"15:11\">201904057691092</ns9:rid></ns9:first></ns9:set></trainOrder></uR></Pport>"
	trainOrderJson = "{\"rid\":\"201904057691092\",\"uid\":\"L91092\",\"trainId\":\"2L81\",\"ssd\":\"2019-04-05\",\"toc\":\"LO\",\"status\":\"P\",\"trainCat\":\"OO\",\"passengerService\":true,\"active\":true,\"cancelReason\":{\"reason\":0},\"lateReason\":{\"reason\":0},\"locations\":[{\"type\":\"OR\",\"tiploc\":\"CLPHMJ1\",\"displaytime\":\"15:02:00\",\"timetable\":{\"time\":\"15:02:00\",\"ptd\":\"15:01\",\"wtd\":\"15:02:00\"},\"planned\":{\"activity\":\"TB\"},\"forecast\":{\"time\":\"15:02:00\",\"departed\":true,\"arr\":null,\"dep\":{\"at\":\"15:02:00\",\"src\":\"TRUST\",\"srcInst\":\"Auto\"},\"pass\":null,\"plat\":{\"plat\":\"1\",\"confirmed\":true,\"source\":\"A\"},\"date\":\"2019-04-05T15:02:36.9040587+01:00\"},\"delay\":0,\"loading\":null},{\"type\":\"PP\",\"tiploc\":\"LTCHMRJ\",\"displaytime\":\"15:04:00\",\"timetable\":{\"time\":\"15:04:00\",\"wtp\":\"15:04:00\"},\"planned\":{},\"forecast\":{\"time\":\"15:04:00\",\"arrived\":true,\"departed\":true,\"arr\":null,\"dep\":null,\"pass\":{\"at\":\"15:04:00\",\"src\":\"TRUST\",\"srcInst\":\"Auto\"},\"plat\":{},\"date\":\"2019-04-05T15:04:07.3216587+01:00\"},\"delay\":0,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"CSEAH\",\"displaytime\":\"15:06:00\",\"timetable\":{\"time\":\"15:06:30\",\"pta\":\"15:05\",\"ptd\":\"15:05\",\"wta\":\"15:06:00\",\"wtd\":\"15:06:30\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:06:00\",\"arrived\":true,\"departed\":true,\"arr\":{\"at\":\"15:05:00\",\"src\":\"TD\"},\"dep\":{\"at\":\"15:06:00\",\"src\":\"TD\"},\"pass\":null,\"plat\":{\"plat\":\"2\",\"confirmed\":true,\"source\":\"A\"},\"trainOrder\":{\"order\":1,\"plat\":\"2\",\"date\":\"2019-04-05T15:05:42.4388587+01:00\"},\"date\":\"2019-04-05T15:07:35.1176587+01:00\"},\"delay\":-30,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"WBRMPTN\",\"displaytime\":\"15:09:00\",\"timetable\":{\"time\":\"15:09:30\",\"pta\":\"15:08\",\"ptd\":\"15:08\",\"wta\":\"15:09:00\",\"wtd\":\"15:09:30\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:09:00\",\"arrived\":true,\"departed\":true,\"arr\":{\"at\":\"15:07:00\",\"src\":\"TD\"},\"dep\":{\"at\":\"15:09:00\",\"src\":\"TD\"},\"pass\":null,\"plat\":{\"plat\":\"4\",\"confirmed\":true,\"source\":\"A\"},\"trainOrder\":{\"order\":1,\"plat\":\"4\",\"date\":\"2019-04-05T15:07:35.3516587+01:00\"},\"date\":\"2019-04-05T15:09:59.1680587+01:00\"},\"delay\":-30,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"KENOLYM\",\"displaytime\":\"15:12:00\",\"timetable\":{\"time\":\"15:12:30\",\"pta\":\"15:11\",\"ptd\":\"15:11\",\"wta\":\"15:12:00\",\"wtd\":\"15:12:30\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:12:00\",\"arrived\":true,\"departed\":true,\"arr\":{\"at\":\"15:11:00\",\"src\":\"TD\"},\"dep\":{\"at\":\"15:12:00\",\"src\":\"TD\"},\"pass\":null,\"plat\":{\"plat\":\"2\",\"confirmed\":true,\"source\":\"A\"},\"trainOrder\":{\"order\":1,\"plat\":\"2\",\"date\":\"2019-04-05T15:09:59.3396587+01:00\"},\"date\":\"2019-04-05T15:12:59.9096587+01:00\"},\"delay\":-30,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"SHPDSB\",\"displaytime\":\"15:15:00\",\"timetable\":{\"time\":\"15:14:30\",\"pta\":\"15:13\",\"ptd\":\"15:13\",\"wta\":\"15:14:00\",\"wtd\":\"15:14:30\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:15:00\",\"arrived\":true,\"departed\":true,\"arr\":{\"at\":\"15:14:00\",\"src\":\"TD\"},\"dep\":{\"at\":\"15:15:00\",\"src\":\"TD\"},\"pass\":null,\"plat\":{\"plat\":\"2\",\"confirmed\":true,\"source\":\"A\"},\"trainOrder\":{\"order\":1,\"plat\":\"2\",\"date\":\"2019-04-05T15:10:49.5716587+01:00\"},\"date\":\"2019-04-05T15:16:05.9240587+01:00\"},\"delay\":30,\"loading\":null},{\"type\":\"OPIP\",\"tiploc\":\"NPLE813\",\"displaytime\":\"15:17:00\",\"timetable\":{\"time\":\"15:16:30\",\"wta\":\"15:16:30\",\"wtd\":\"15:16:30\"},\"planned\":{\"activity\":\"OP\"},\"forecast\":{\"time\":\"15:17:00\",\"arrived\":true,\"departed\":true,\"arr\":{\"at\":\"15:15:00\",\"src\":\"TD\"},\"dep\":{\"at\":\"15:17:00\",\"src\":\"TRUST\",\"srcInst\":\"Auto\"},\"pass\":null,\"plat\":{},\"date\":\"2019-04-05T15:17:25.9988587+01:00\"},\"delay\":30,\"loading\":null},{\"type\":\"PP\",\"tiploc\":\"NPOLEJ\",\"displaytime\":\"15:18:00\",\"timetable\":{\"time\":\"15:17:30\",\"wtp\":\"15:17:30\"},\"planned\":{},\"forecast\":{\"time\":\"15:18:00\",\"arr\":null,\"dep\":null,\"pass\":{\"et\":\"15:18:00\",\"src\":\"Darwin\"},\"plat\":{},\"date\":\"2019-04-05T15:17:01.3664587+01:00\"},\"delay\":30,\"loading\":null},{\"type\":\"PP\",\"tiploc\":\"MTRBDGJ\",\"displaytime\":\"15:18:00\",\"timetable\":{\"time\":\"15:18:00\",\"wtp\":\"15:18:00\"},\"planned\":{},\"forecast\":{\"time\":\"15:18:00\",\"arrived\":true,\"departed\":true,\"arr\":null,\"dep\":null,\"pass\":{\"at\":\"15:18:00\",\"src\":\"TRUST\",\"srcInst\":\"Auto\"},\"plat\":{},\"date\":\"2019-04-05T15:18:47.1968587+01:00\"},\"delay\":0,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"WLSDJHL\",\"displaytime\":\"15:22:00\",\"timetable\":{\"time\":\"15:23:00\",\"pta\":\"15:21\",\"ptd\":\"15:22\",\"wta\":\"15:21:30\",\"wtd\":\"15:23:00\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:22:00\",\"arrived\":true,\"departed\":true,\"arr\":{\"at\":\"15:21:00\",\"src\":\"TD\"},\"dep\":{\"at\":\"15:22:00\",\"src\":\"TD\"},\"pass\":null,\"plat\":{\"plat\":\"4\",\"confirmed\":true,\"source\":\"A\"},\"date\":\"2019-04-05T15:23:03.7076587+01:00\"},\"delay\":-60,\"loading\":null},{\"type\":\"PP\",\"tiploc\":\"KENSLGJ\",\"displaytime\":\"15:24:00\",\"timetable\":{\"time\":\"15:24:00\",\"wtp\":\"15:24:00\"},\"planned\":{},\"forecast\":{\"time\":\"15:24:00\",\"arrived\":true,\"departed\":true,\"arr\":null,\"dep\":null,\"pass\":{\"at\":\"15:24:00\",\"src\":\"TRUST\",\"srcInst\":\"Auto\"},\"plat\":{},\"date\":\"2019-04-05T15:24:21.1616587+01:00\"},\"delay\":0,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"KENR\",\"displaytime\":\"15:26:00\",\"timetable\":{\"time\":\"15:25:30\",\"pta\":\"15:24\",\"ptd\":\"15:24\",\"wta\":\"15:25:00\",\"wtd\":\"15:25:30\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:26:00\",\"arrived\":true,\"departed\":true,\"arr\":{\"at\":\"15:24:00\",\"src\":\"TD\"},\"dep\":{\"at\":\"15:26:00\",\"src\":\"TD\"},\"pass\":null,\"plat\":{\"plat\":\"2\",\"confirmed\":true,\"source\":\"A\"},\"trainOrder\":{\"order\":1,\"plat\":\"2\",\"date\":\"2019-04-05T15:24:51.5263383+01:00\"},\"date\":\"2019-04-05T15:26:31.3008719+01:00\"},\"delay\":30,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"BRBYPK\",\"displaytime\":\"15:28:00\",\"timetable\":{\"time\":\"15:27:30\",\"pta\":\"15:26\",\"ptd\":\"15:26\",\"wta\":\"15:27:00\",\"wtd\":\"15:27:30\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:28:00\",\"arrived\":true,\"departed\":true,\"arr\":{\"at\":\"15:26:00\",\"src\":\"TD\"},\"dep\":{\"at\":\"15:28:00\",\"src\":\"TD\"},\"pass\":null,\"plat\":{\"plat\":\"2\",\"confirmed\":true,\"source\":\"A\"},\"date\":\"2019-04-05T15:28:39.4129959+01:00\"},\"delay\":30,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"BRBY\",\"displaytime\":\"15:30:00\",\"timetable\":{\"time\":\"15:29:00\",\"pta\":\"15:28\",\"ptd\":\"15:28\",\"wta\":\"15:28:30\",\"wtd\":\"15:29:00\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:30:00\",\"arrived\":true,\"departed\":true,\"arr\":{\"at\":\"15:29:00\",\"src\":\"TD\"},\"dep\":{\"at\":\"15:30:00\",\"src\":\"TD\"},\"pass\":null,\"plat\":{\"plat\":\"2\",\"confirmed\":true,\"source\":\"A\"},\"trainOrder\":{\"order\":1,\"plat\":\"2\",\"date\":\"2019-04-05T15:26:31.4881247+01:00\"},\"date\":\"2019-04-05T15:30:41.4103567+01:00\"},\"delay\":60,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"WHMDSTD\",\"displaytime\":\"15:32:00\",\"timetable\":{\"time\":\"15:31:00\",\"pta\":\"15:30\",\"ptd\":\"15:30\",\"wta\":\"15:30:30\",\"wtd\":\"15:31:00\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:32:00\",\"arrived\":true,\"departed\":true,\"arr\":{\"at\":\"15:31:00\",\"src\":\"TD\"},\"dep\":{\"at\":\"15:32:00\",\"src\":\"TD\"},\"pass\":null,\"plat\":{\"plat\":\"2\",\"confirmed\":true,\"source\":\"A\"},\"trainOrder\":{\"order\":1,\"plat\":\"2\",\"date\":\"2019-04-05T15:28:39.6314575+01:00\"},\"date\":\"2019-04-05T15:32:29.3734287+01:00\"},\"delay\":60,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"FNCHLYR\",\"displaytime\":\"15:33:00\",\"timetable\":{\"time\":\"15:32:30\",\"pta\":\"15:31\",\"ptd\":\"15:31\",\"wta\":\"15:32:00\",\"wtd\":\"15:32:30\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:33:00\",\"arrived\":true,\"departed\":true,\"arr\":{\"at\":\"15:33:00\",\"src\":\"TD\"},\"dep\":{\"at\":\"15:33:00\",\"src\":\"TD\"},\"pass\":null,\"plat\":{\"plat\":\"2\",\"confirmed\":true,\"source\":\"A\"},\"trainOrder\":{\"order\":1,\"plat\":\"2\",\"date\":\"2019-04-05T15:30:41.5663727+01:00\"},\"date\":\"2019-04-05T15:33:57.4444607+01:00\"},\"delay\":30,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"HMPSTDH\",\"displaytime\":\"15:35:00\",\"timetable\":{\"time\":\"15:35:00\",\"pta\":\"15:34\",\"ptd\":\"15:34\",\"wta\":\"15:34:30\",\"wtd\":\"15:35:00\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:35:00\",\"arrived\":true,\"arr\":{\"at\":\"15:34:00\",\"src\":\"TD\"},\"dep\":{\"et\":\"15:35:00\",\"src\":\"Darwin\"},\"pass\":null,\"plat\":{\"plat\":\"2\",\"confirmed\":true,\"source\":\"A\"},\"trainOrder\":{\"order\":1,\"plat\":\"2\",\"date\":\"2019-04-05T15:32:30.1067039+01:00\"},\"date\":\"2019-04-05T15:34:47.4618493+01:00\"},\"delay\":0,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"GOSPLOK\",\"displaytime\":\"15:37:00\",\"timetable\":{\"time\":\"15:37:30\",\"pta\":\"15:36\",\"ptd\":\"15:36\",\"wta\":\"15:37:00\",\"wtd\":\"15:37:30\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:37:00\",\"arr\":{\"et\":\"15:37:00\",\"src\":\"Darwin\"},\"dep\":{\"et\":\"15:37:00\",\"src\":\"Darwin\"},\"pass\":null,\"plat\":{\"plat\":\"2\",\"confirmed\":true,\"source\":\"A\"},\"trainOrder\":{\"order\":1,\"plat\":\"2\",\"date\":\"2019-04-05T15:33:59.4570671+01:00\"},\"date\":\"2019-04-05T15:34:47.4618493+01:00\"},\"delay\":-30,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"KNTSHTW\",\"displaytime\":\"15:39:00\",\"timetable\":{\"time\":\"15:39:30\",\"pta\":\"15:38\",\"ptd\":\"15:38\",\"wta\":\"15:39:00\",\"wtd\":\"15:39:30\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:39:00\",\"arr\":{\"et\":\"15:39:00\",\"src\":\"Darwin\"},\"dep\":{\"et\":\"15:39:00\",\"src\":\"Darwin\"},\"pass\":null,\"plat\":{\"plat\":\"2\"},\"date\":\"2019-04-05T15:34:47.4618493+01:00\"},\"delay\":-30,\"loading\":null},{\"type\":\"PP\",\"tiploc\":\"CMDNRDJ\",\"displaytime\":\"15:42:00\",\"timetable\":{\"time\":\"15:42:00\",\"wtp\":\"15:42:00\"},\"planned\":{},\"forecast\":{\"time\":\"15:42:00\",\"arr\":null,\"dep\":null,\"pass\":{\"et\":\"15:42:00\",\"src\":\"Darwin\"},\"plat\":{},\"date\":\"2019-04-05T15:34:47.4618493+01:00\"},\"delay\":0,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"CMDNRD\",\"displaytime\":\"15:43:00\",\"timetable\":{\"time\":\"15:43:00\",\"pta\":\"15:43\",\"ptd\":\"15:43\",\"wta\":\"15:42:30\",\"wtd\":\"15:43:00\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:43:00\",\"arr\":{\"et\":\"15:43:00\",\"wet\":\"15:42:00\",\"src\":\"Darwin\"},\"dep\":{\"et\":\"15:43:00\",\"src\":\"Darwin\"},\"pass\":null,\"plat\":{\"plat\":\"2\"},\"date\":\"2019-04-05T15:34:47.4618493+01:00\"},\"delay\":0,\"loading\":null},{\"type\":\"PP\",\"tiploc\":\"CMDNREJ\",\"displaytime\":\"15:44:00\",\"timetable\":{\"time\":\"15:44:00\",\"wtp\":\"15:44:00\"},\"planned\":{},\"forecast\":{\"time\":\"15:44:00\",\"arr\":null,\"dep\":null,\"pass\":{\"et\":\"15:44:00\",\"src\":\"Darwin\"},\"plat\":{},\"date\":\"2019-04-05T15:34:47.4618493+01:00\"},\"delay\":0,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"CLDNNRB\",\"displaytime\":\"15:46:00\",\"timetable\":{\"time\":\"15:46:00\",\"pta\":\"15:46\",\"ptd\":\"15:46\",\"wta\":\"15:45:30\",\"wtd\":\"15:46:00\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:46:00\",\"arr\":{\"et\":\"15:46:00\",\"wet\":\"15:45:00\",\"src\":\"Darwin\"},\"dep\":{\"et\":\"15:46:00\",\"src\":\"Darwin\"},\"pass\":null,\"plat\":{\"plat\":\"3\"},\"date\":\"2019-04-05T15:34:47.4618493+01:00\"},\"delay\":0,\"loading\":null},{\"type\":\"PP\",\"tiploc\":\"WSBRNRJ\",\"displaytime\":\"15:47:00\",\"timetable\":{\"time\":\"15:47:00\",\"wtp\":\"15:47:00\"},\"planned\":{},\"forecast\":{\"time\":\"15:47:00\",\"arr\":null,\"dep\":null,\"pass\":{\"et\":\"15:47:00\",\"src\":\"Darwin\"},\"plat\":{},\"date\":\"2019-04-05T15:34:47.4618493+01:00\"},\"delay\":0,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"HIGHBYA\",\"displaytime\":\"15:48:00\",\"timetable\":{\"time\":\"15:48:30\",\"pta\":\"15:48\",\"ptd\":\"15:48\",\"wta\":\"15:47:30\",\"wtd\":\"15:48:30\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:48:00\",\"arr\":{\"et\":\"15:48:00\",\"wet\":\"15:47:00\",\"src\":\"Darwin\"},\"dep\":{\"et\":\"15:48:00\",\"src\":\"Darwin\"},\"pass\":null,\"plat\":{\"plat\":\"8\"},\"date\":\"2019-04-05T15:34:47.4618493+01:00\"},\"delay\":-30,\"loading\":null},{\"type\":\"PP\",\"tiploc\":\"CNNBYWJ\",\"displaytime\":\"15:49:00\",\"timetable\":{\"time\":\"15:49:30\",\"wtp\":\"15:49:30\"},\"planned\":{},\"forecast\":{\"time\":\"15:49:00\",\"arr\":null,\"dep\":null,\"pass\":{\"et\":\"15:49:00\",\"src\":\"Darwin\"},\"plat\":{},\"date\":\"2019-04-05T15:34:47.4618493+01:00\"},\"delay\":-30,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"CNNB\",\"displaytime\":\"15:50:00\",\"timetable\":{\"time\":\"15:50:30\",\"pta\":\"15:50\",\"ptd\":\"15:50\",\"wta\":\"15:50:00\",\"wtd\":\"15:50:30\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:50:00\",\"arr\":{\"et\":\"15:50:00\",\"wet\":\"15:49:00\",\"src\":\"Darwin\"},\"dep\":{\"et\":\"15:50:00\",\"src\":\"Darwin\"},\"pass\":null,\"plat\":{\"plat\":\"4\"},\"date\":\"2019-04-05T15:34:47.4618493+01:00\"},\"delay\":-30,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"DALSKLD\",\"displaytime\":\"15:52:00\",\"timetable\":{\"time\":\"15:52:30\",\"pta\":\"15:52\",\"ptd\":\"15:52\",\"wta\":\"15:52:00\",\"wtd\":\"15:52:30\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:52:00\",\"arr\":{\"et\":\"15:52:00\",\"wet\":\"15:51:00\",\"src\":\"Darwin\"},\"dep\":{\"et\":\"15:52:00\",\"src\":\"Darwin\"},\"pass\":null,\"plat\":{\"plat\":\"2\"},\"date\":\"2019-04-05T15:34:47.4618493+01:00\"},\"delay\":-30,\"loading\":null},{\"type\":\"PP\",\"tiploc\":\"NAVRDJ\",\"displaytime\":\"15:53:00\",\"timetable\":{\"time\":\"15:53:30\",\"wtp\":\"15:53:30\"},\"planned\":{},\"forecast\":{\"time\":\"15:53:00\",\"arr\":null,\"dep\":null,\"pass\":{\"et\":\"15:53:00\",\"src\":\"Darwin\"},\"plat\":{},\"date\":\"2019-04-05T15:34:47.4618493+01:00\"},\"delay\":-30,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"HACKNYC\",\"displaytime\":\"15:55:00\",\"timetable\":{\"time\":\"15:55:00\",\"pta\":\"15:55\",\"ptd\":\"15:55\",\"wta\":\"15:54:30\",\"wtd\":\"15:55:00\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:55:00\",\"arr\":{\"et\":\"15:55:00\",\"wet\":\"15:53:00\",\"src\":\"Darwin\"},\"dep\":{\"et\":\"15:55:00\",\"src\":\"Darwin\"},\"pass\":null,\"plat\":{\"plat\":\"2\"},\"date\":\"2019-04-05T15:34:47.4618493+01:00\"},\"delay\":0,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"HOMRTON\",\"displaytime\":\"15:57:00\",\"timetable\":{\"time\":\"15:57:00\",\"pta\":\"15:57\",\"ptd\":\"15:57\",\"wta\":\"15:56:30\",\"wtd\":\"15:57:00\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:57:00\",\"arr\":{\"et\":\"15:57:00\",\"wet\":\"15:56:00\",\"src\":\"Darwin\"},\"dep\":{\"et\":\"15:57:00\",\"src\":\"Darwin\"},\"pass\":null,\"plat\":{\"plat\":\"2\"},\"date\":\"2019-04-05T15:33:57.4444607+01:00\"},\"delay\":0,\"loading\":null},{\"type\":\"IP\",\"tiploc\":\"HACKNYW\",\"displaytime\":\"15:59:00\",\"timetable\":{\"time\":\"15:59:30\",\"pta\":\"15:59\",\"ptd\":\"15:59\",\"wta\":\"15:59:00\",\"wtd\":\"15:59:30\"},\"planned\":{\"activity\":\"T \"},\"forecast\":{\"time\":\"15:59:00\",\"arr\":{\"et\":\"15:59:00\",\"src\":\"Darwin\"},\"dep\":{\"et\":\"15:59:00\",\"src\":\"Darwin\"},\"pass\":null,\"plat\":{\"plat\":\"2\"},\"date\":\"2019-04-05T15:33:57.4444607+01:00\"},\"delay\":-30,\"loading\":null},{\"type\":\"PP\",\"tiploc\":\"LEAJ\",\"displaytime\":\"16:00:00\",\"timetable\":{\"time\":\"16:00:30\",\"wtp\":\"16:00:30\"},\"planned\":{},\"forecast\":{\"time\":\"16:00:00\",\"arr\":null,\"dep\":null,\"pass\":{\"et\":\"16:00:00\",\"src\":\"Darwin\"},\"plat\":{},\"date\":\"2019-04-05T15:33:57.4444607+01:00\"},\"delay\":-30,\"loading\":null},{\"type\":\"PP\",\"tiploc\":\"CHNELSJ\",\"displaytime\":\"16:01:00\",\"timetable\":{\"time\":\"16:01:00\",\"wtp\":\"16:01:00\"},\"planned\":{},\"forecast\":{\"time\":\"16:01:00\",\"arr\":null,\"dep\":null,\"pass\":null,\"plat\":{},\"date\":\"0001-01-01T00:00:00Z\"},\"delay\":0,\"loading\":null},{\"type\":\"DT\",\"tiploc\":\"STFD\",\"displaytime\":\"16:05:00\",\"timetable\":{\"time\":\"16:03:00\",\"pta\":\"16:05\",\"wta\":\"16:03:00\"},\"planned\":{\"activity\":\"TF\"},\"forecast\":{\"time\":\"16:05:00\",\"arr\":null,\"dep\":null,\"pass\":null,\"plat\":{},\"date\":\"0001-01-01T00:00:00Z\"},\"delay\":120,\"loading\":null}],\"originLocation\":{\"type\":\"OR\",\"tiploc\":\"CLPHMJ1\",\"displaytime\":\"15:02:00\",\"timetable\":{\"time\":\"15:02:00\",\"ptd\":\"15:01\",\"wtd\":\"15:02:00\"},\"planned\":{\"activity\":\"TB\"},\"forecast\":{\"time\":\"15:02:00\",\"departed\":true,\"arr\":null,\"dep\":{\"at\":\"15:02:00\",\"src\":\"TRUST\",\"srcInst\":\"Auto\"},\"pass\":null,\"plat\":{\"plat\":\"1\",\"confirmed\":true,\"source\":\"A\"},\"date\":\"2019-04-05T15:02:36.9040587+01:00\"},\"delay\":0,\"loading\":null},\"destinationLocation\":{\"type\":\"DT\",\"tiploc\":\"STFD\",\"displaytime\":\"16:05:00\",\"timetable\":{\"time\":\"16:03:00\",\"pta\":\"16:05\",\"wta\":\"16:03:00\"},\"planned\":{\"activity\":\"TF\"},\"forecast\":{\"time\":\"16:05:00\",\"arr\":null,\"dep\":null,\"pass\":null,\"plat\":{},\"date\":\"0001-01-01T00:00:00Z\"},\"delay\":120,\"loading\":null},\"association\":null,\"date\":\"2019-04-05T15:34:47.4618493+01:00\"}"
)

// Test_TrainOrder_XML_Parse checks that we can parse normal xml
func Test_TrainOrder_XML_Parse(t *testing.T) {

	type pport struct {
		TS time.Time `xml:"ts,attr"`
		UR struct {
			Order trainOrderWrapper `xml:"trainOrder"`
		} `xml:"uR"`
	}
	p := &pport{}

	r := bytes.NewReader([]byte(trainOrderXML))
	decoder := xml.NewDecoder(r)
	err := decoder.Decode(p)

	if err != nil {
		t.Errorf("Failed to parse xml: %s", err)
	}
}

// Test_Issue6_TrainOrder_XML_Parse looks to ensure that we can decode an xml message correctly.
// issue6Xml contains a live instance of this where the xml parser failed to process this line
func Test_Issue6_TrainOrder_XML_Parse(t *testing.T) {

	type pport struct {
		TS time.Time `xml:"ts,attr"`
		UR struct {
			Order trainOrderWrapper `xml:"trainOrder"`
		} `xml:"uR"`
	}
	p := &pport{}

	r := bytes.NewReader([]byte(issue6Xml))
	decoder := xml.NewDecoder(r)
	err := decoder.Decode(p)

	if err != nil {
		t.Errorf("Failed to parse xml")
	}

	obj := p.UR.Order

	if obj.Tiploc != "BRGEND" {
		t.Errorf("Tiploc \"%s\" expected \"BRGEND\"", obj.Tiploc)
	}

	if obj.CRS != "BGN" {
		t.Errorf("Crs\"%s\" expected \"BGN\"", obj.CRS)
	}

	if obj.Platform != "1" {
		t.Errorf("Platform \"%s\" expected \"1\"", obj.Platform)
	}

	if obj.Set == nil {
		t.Error("Set is nil")
	} else if obj.Set.First == nil {
		t.Error("Set.First is nil")
	} else {
		f := obj.Set.First

		if f.RID.RID != "" {
			t.Errorf("RID \"%s\" expected \"\"", f.RID.RID)
		}

		// TODO test f.Times

		if f.TrainId != "6B09" {
			t.Errorf("TrainId \"%s\" expected \"6B09\"", f.TrainId)
		}
	}

	if obj.Set.Second != nil {
		t.Error("Set.Second is not nil")
	}

	if obj.Set.Third != nil {
		t.Error("Set.Third is not nil")
	}
}

// Test_TrainOrder_Apply tests that we can apply using constructed data
func Test_TrainOrder_Apply(t *testing.T) {

	tpl := "SHPDSB"
	plat := "1"

	loc := &Location{
		Tiploc: tpl,
		Times: util.CircularTimes{
			Pta: util.NewPublicTime("13:11"),
			Ptd: util.NewPublicTime("13:11"),
			Wta: util.NewWorkingTime("13:11"),
			Wtd: util.NewWorkingTime("13:11:30"),
		},
	}

	sched := &Schedule{
		Locations: []*Location{
			&Location{
				Tiploc: "MSTONEE",
				Times: util.CircularTimes{
					Pta: util.NewPublicTime("13:00"),
					Ptd: util.NewPublicTime("13:00"),
					Wta: util.NewWorkingTime("13:00"),
					Wtd: util.NewWorkingTime("13:00"),
				},
			},
			loc,
			&Location{
				Tiploc: "KENOLYM",
				Times: util.CircularTimes{
					Pta: util.NewPublicTime("13:30"),
					Ptd: util.NewPublicTime("13:30"),
					Wta: util.NewWorkingTime("13:30"),
					Wtd: util.NewWorkingTime("13:30"),
				},
			},
		},
	}

	tod := &trainOrderItem{
		RID: TrainOrderRID{
			RID: "201904057692589",
			Times: util.CircularTimes{
				Pta: util.NewPublicTime("13:11"),
				Ptd: util.NewPublicTime("13:11"),
				Wta: util.NewWorkingTime("13:11"),
				Wtd: util.NewWorkingTime("13:11:30"),
			},
		},
	}

	to := &trainOrderWrapper{
		Tiploc:   tpl,
		Platform: plat,
		Set: &trainOrderData{
			First: tod,
		},
	}

	ts := time.Now()

	l := tod.apply(to, 1, ts, sched)
	if !l {
		t.Error("Failed to find location in schedule")
		return
	}

	if loc.Forecast.TrainOrder == nil {
		t.Error("TrainOrder not applied to forecast")
	}

	b, _ := sched.Bytes()
	log.Println(string(b))
}

// Test_TrainOrder_ApplyXml tests we can apply using live xml
func Test_TrainOrder_ApplyXml(t *testing.T) {

	type pport struct {
		TS time.Time `xml:"ts,attr"`
		UR struct {
			Order trainOrderWrapper `xml:"trainOrder"`
		} `xml:"uR"`
	}
	p := &pport{}

	r := bytes.NewReader([]byte(trainOrderXML))
	decoder := xml.NewDecoder(r)
	err := decoder.Decode(p)

	if err != nil {
		t.Errorf("Failed to parse xml: %s", err)
	}

	sched := ScheduleFromBytes([]byte(trainOrderJson))

	ts := time.Now()

	l := p.UR.Order.Set.First.apply(&p.UR.Order, 1, ts, sched)
	if !l {
		t.Error("Failed to find location in schedule")
		return
	}

	tpl := "KENOLYM"
	found := false
	for _, loc := range sched.Locations {
		if loc.Tiploc == tpl {
			found = true
			if loc.Forecast.TrainOrder == nil {
				t.Error("TrainOrder not applied to forecast")
			}
		}
	}

	if !found {
		t.Error("Failed to find location in response")
	}

}
