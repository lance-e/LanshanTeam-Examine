package consts

const (
	LackParams = 1000 + iota
	UserNotFound
	UserAlreadyExist
	RegisterFailed
	RegisterSuccess
	LoginPasswordWrong
	LoginSuccess
	LoginFailed
	GenerateTokenFailed
	TokenInvalid
	ServeUnavailable
	SendCodeSuccess
	SendCodeFailed
	CodeWrong
	PhoneNumberUnavailable
	GetUserAllInformationSuccess
	GetUserAllInformationFailed
)
const (
	AddFriendRequestSuccess = 1015 + iota
	CreateFriendSuccess
	StartGameFailed
	StartGameSuccess
	NotFoundTargetGameRoom
	ReadyToPlayGameFailed
	ReadyToPlayGameSuccess
	GameOver
)
