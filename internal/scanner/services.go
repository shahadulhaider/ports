package scanner

var serviceNames = map[int]string{
	21:    "ftp",
	22:    "ssh",
	25:    "smtp",
	53:    "dns",
	80:    "http",
	110:   "pop3",
	143:   "imap",
	443:   "https",
	993:   "imaps",
	995:   "pop3s",
	1433:  "mssql",
	1521:  "oracle",
	3000:  "dev",
	3306:  "mysql",
	4200:  "ng-serve",
	5000:  "flask",
	5432:  "postgres",
	5672:  "amqp",
	5900:  "vnc",
	6379:  "redis",
	8000:  "dev",
	8080:  "http-alt",
	8443:  "https-alt",
	8888:  "jupyter",
	9090:  "prometheus",
	9200:  "elastic",
	11211: "memcached",
	27017: "mongodb",
}

func ServiceName(port int) string {
	return serviceNames[port]
}
