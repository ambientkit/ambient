module github.com/josephspurrier/ambient

go 1.16

require (
	cloud.google.com/go/storage v1.18.2
	github.com/Azure/azure-storage-blob-go v0.14.0
	github.com/alexedwards/scs/v2 v2.5.0
	github.com/aws/aws-sdk-go-v2 v1.13.0
	github.com/aws/aws-sdk-go-v2/config v1.13.0
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.9.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.24.0
	github.com/c-bata/go-prompt v0.2.6
	github.com/dgryski/dgoogauth v0.0.0-20190221195224-5a805980a5f3
	github.com/mattn/go-colorable v0.1.12
	github.com/microcosm-cc/bluemonday v1.0.17
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/sirupsen/logrus v1.8.1
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/ssor/bom v0.0.0-20170718123548-6386211fdfcf // indirect
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.20.0
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce
	jaytaylor.com/html2text v0.0.0-20211105163654-bc68cce691ba
)

replace github.com/c-bata/go-prompt => github.com/josephspurrier/go-prompt v0.2.7-0.20220117021137-67747317bc02
