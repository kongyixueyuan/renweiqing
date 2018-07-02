package BLC

type TXOutput struct {
	Value int64 // 金额/面值
	ScriptPubKey string  //用户名
}
// 解锁
func (txOutput *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {
	return txOutput.ScriptPubKey == address
}