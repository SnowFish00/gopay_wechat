package model_cfg

import (
	"github.com/go-pay/gopay/wechat/v3"
)

type ReflectByConfig interface {
	NewClientV3Engine() *wechat.ClientV3
}
