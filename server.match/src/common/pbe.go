package common

type PbeRight struct {
	Id    uint64
	Code  string
	Right uint32
}

type ReqGetPbeRight struct {
	UID uint64
}

type RetGetPbeRight struct {
	Right uint32
}

type ReqPbeReject struct {
	UID    uint64
	Reject uint32
}

type RetPbeReject struct {
}

type ReqPbeSureyAnswer struct {
	UID    uint64
	Answer string
}

type RetPbeSureyAnswer struct {
	Code string
}

type ReqPbeRightCode struct {
	UID uint64
}

type RetPbeRightCode struct {
	Right uint32
	Code  string
}

type ReqPbeAdvice struct {
	UID     uint64
	Advice  string
	PUrls   []string
	Account string
	Type    string
}

type RetPbeAdvice struct {
}

type ReqCheckPbeCode struct {
	Dev  string
	Code string
}

type RetCheckPbeCode struct {
	ErrCode int
}
