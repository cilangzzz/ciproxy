package ciproxy

import "time"

// ProxyMethod constant proxyMethod
const (
	HttpProxy             = "HttpProxy"
	HttpsProxy            = "HttpsProxy"
	HttpsSniffProxy       = "HttpsSniffProxy"
	HttpsSniffDetailProxy = "HttpsSniffDetailProxy"
	WebsocketProxy        = "WebsocketProxy"
	TcpNormalProxy        = "TcpNormal"
	TcpTunnelProxy        = "TcpTunnel"
	PortProxy             = "PortProxy"
	DefaultProxy          = "All"
)

// Proxy Config
const (
	ProxyVersion = "v0.0.0"
	ProxyMode    = "Debug"

	// DefaultIp DefaultPort defaultServerConfig
	DefaultIp   = "127.0.0.1"
	DefaultPort = ""

	ProxyOrganization = "www.cilang.buzz"
)

// connectConfig Connect Config Constant
const (
	DefaultConnectProtocol = "Tcp"
	DefaultOutTime         = 10 * time.Second
)

// DefaultCert 默认证书
var DefaultCert = []byte("-----BEGIN CERTIFICATE-----\nMIIDrzCCApcCFCO3T1AFJW49FIgStat27COPWomaMA0GCSqGSIb3DQEBCwUAMIGT\nMQswCQYDVQQGEwJDTjEQMA4GA1UECAwHQ2lwcm94eTESMBAGA1UEBwwJR3Vhbmda\naG91MRAwDgYDVQQKDAdDaXByb3h5MRUwEwYDVQQLDAxDaXByb3h5SHR0cHMxEDAO\nBgNVBAMMB0NpcHJveHkxIzAhBgkqhkiG9w0BCQEWFGNpbGFuZ3VzZXJAZ21haWwu\nY29tMB4XDTIzMTIyNjA4MDkwMVoXDTI0MTIyNTA4MDkwMVowgZMxCzAJBgNVBAYT\nAkNOMRAwDgYDVQQIDAdDaXByb3h5MRIwEAYDVQQHDAlHdWFuZ1pob3UxEDAOBgNV\nBAoMB0NpcHJveHkxFTATBgNVBAsMDENpcHJveHlIdHRwczEQMA4GA1UEAwwHQ2lw\ncm94eTEjMCEGCSqGSIb3DQEJARYUY2lsYW5ndXNlckBnbWFpbC5jb20wggEiMA0G\nCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDb52u1P+BlDRCYxBLfiywnNWewD0sN\nEa3CXZyWSzhkCOgaMO1AwXHljmh3DwhI2zEaD1BFkOJPJd8Lw+QaFoSEe+wS5uqP\nOTylK0OUfGbku1SV48mkhcirA5hs3zDw65pRm8H7UOTIlqtcT0LEKXlTibx2iSNC\nZ+u0U2QUfAtWgWVrL75ssOKpz4CmpV31xbLrDpL3itcZVXj+w5qtL77xg9r5mAca\nLft88MjS7XffNW68HsIJ0zXoaLLvr9HXdS6kr51CXYJc0RZhgljPQHLd/i/GIBzh\nu3t72yqTDQ9EoiDJ/BVt2mwNbM7bQJw+TLlhS4cWKqrU1sXoJ0iMVd3hAgMBAAEw\nDQYJKoZIhvcNAQELBQADggEBAMdWET19NpQ6ZW4RmWIfat4NyOvOJvEKSPF/gNLs\n15mzlpESSL2EVTM+BOEwwsKLH/STcFR84Pb+KdTimAW8Tl6kD5xUdbxgGZKvxIeq\ngYrzkqLS8SznyDErIbNux8hdGUnvlbLb8e1w4GumOKUqP2Zp0vbQ7B9QkdQ4c6Ly\nHbhsAfj4ULCeaJyqFXf1Ar+A2XDxR2C/aXM6fKFwyD+0xciQTXNbA/VH2WSokESJ\nM3fE8stCi3chChqqbAaMSsu1a6hQzqfPM/KdowvjED/M3DPliJ83SVxGeyblV1/j\n877WLawMwlMEUn35+8SQBVTyEM4EEnG8wH/35kdrIxSVqnQ=\n-----END CERTIFICATE-----\n")

// DefaultCertKey 默认私钥
var DefaultCertKey = []byte("-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDb52u1P+BlDRCY\nxBLfiywnNWewD0sNEa3CXZyWSzhkCOgaMO1AwXHljmh3DwhI2zEaD1BFkOJPJd8L\nw+QaFoSEe+wS5uqPOTylK0OUfGbku1SV48mkhcirA5hs3zDw65pRm8H7UOTIlqtc\nT0LEKXlTibx2iSNCZ+u0U2QUfAtWgWVrL75ssOKpz4CmpV31xbLrDpL3itcZVXj+\nw5qtL77xg9r5mAcaLft88MjS7XffNW68HsIJ0zXoaLLvr9HXdS6kr51CXYJc0RZh\ngljPQHLd/i/GIBzhu3t72yqTDQ9EoiDJ/BVt2mwNbM7bQJw+TLlhS4cWKqrU1sXo\nJ0iMVd3hAgMBAAECggEAEdMcIupVLDu3EnrwmktNcYKQ4tXkHBOu+I539W+OPmWO\nbFkQHuvPNjx116vrQv4paIn0B4xRgix5FbRzyXDCN+IhMslu59vlQFbvOM3HCIrt\n6OLSFgWvsogtecQZoBDcE5NcuBwVf1p4wP6K7ejJG24sGmIBVQ4K+zUea2tzZWtP\nKGqyGC2uDK0wmtjkmLnhIIyXTB0ZaxiDnskVmFbqeCZlY5784cAFxCLkl72g1ask\nu5ht7a+c83lg3Yo77ocuF+LQIp5NRIwd4BivdaoJjnbOWeImT5soGqt3JGzDLzLT\nzFYerDIgwjlBJS+xJOYzfk4m/fye5PAtiRskTY4vNQKBgQD2z34x8U1t7wbLWy6B\nldZ30SFWFUwT8AmyTOn6ZghS1rxZ260FxgPcTI9/YPbZf/jE8l0DuSny8fMThdRy\nSoy8vL7goHmZwXiT4WtveAIZJwH2WMfSQjl7cBDLEfM25eCN783Qt1nNqpZehscv\nW10zgHf0RhgTcxZpScxUHPrYhwKBgQDkF3bxdS7IK2eFuoBN3RA/rZP+3dmuJnmq\n6aP7GPPqVTPv295a3GU1Pj8dAHj+4lshJoHHReXsLR786Gu9gV2di2ibI14xn5ZY\n6TIkORJDrU9wkpKGlZsreXcH8lXD8GX2eW+/9eSMtkPkPx1CN/dl0QSC84VyHOl3\nEkUkpnF4VwKBgQDBsA4h1XNlTYqwdgsmKNeZSeZ2bto4X0zMvy1zkzT/BYPkpM/A\n0yfeb7rBRPATuikZLfYu4NX50URoUsNpOfX+e8Tz9RvYvQsKSoIMhUpoQMN6dyvB\nZrVDmKulIZ4TvA0gdku3etwf2bqNzglssa+Ppkb8zTFBujShRgfzDpcQ6wKBgQDC\nlFHfwjvaf0ydBHEX+5I1AHrUXaWpryDz8MT3HF5Ydb8My+BwWrvsI+Hrd8/AgJGF\nQjhMKxDt3cAluJDQ5V9bWuYuEon0g1PbHXSs+hLesSanviJJta7d84zbtpv7v0T8\ncrQmajpC3+oi+MSZDO18akcS/3PD2W5BKdDaZzM9mQKBgHpkLEww2/NthFXg5F0E\nPS23WmBtz5921tVdg1MlAXtZzs8ordT0+w8enGJHglWdBVzrF4qdcTN2tjnlu2e/\nNhGyJiJZ65+dINTibcV3712ChCwvjbf0rTGnpbXaXQHqWmZaC36DchCkdfbJMUCY\nsXnkBxS7/+CrzLBtcmLhaKVJ\n-----END PRIVATE KEY-----\n")
