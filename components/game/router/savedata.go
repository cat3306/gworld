package router

func SaveData() {
	PlayerManager.SaveData()
}

func ClientDisConnect(clientId string) {
	PlayerManager.Remove(clientId)
}
