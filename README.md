# checkup
```
{
	"checkup":{
		"hostName":"自定义名称", //自定义运行设备的名称
		"interval":10, //检测周期
		"pingTimeout":300, //检测单个节点时超时时长
		"pingRetry":4, //ping检测重复次数
		"fastPingMode": false,//是否开启false模式
		"failureRate":0.8,//ipRange内的地址达到80%都不通时就会执行add iptables
		"postUrl": "http://192.168.99.16:8000", //发送微信报警的API
		"ipRange":[ //要检测的节点IP列表,可以上多个哦
			"192.168.99.233",
			"114.114.114.114",
			"218.107.49.162",
			"121.33.191.157",
			"221.5.88.88"
 		]

 	}
}
```
