package loop_err

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

// CusCode is a custom error code type.
// It consists of 7 digits: the first 3 represent the HTTP standard status code,
// and the last 4 are our custom error codes.
// The last 4 digits start from 0000 and increase sequentially (skipping numbers is prohibited).
type CusCode int

const (
	// 200 status code from here
	OK        = 200_00_00  // 成功
	Created   = 201_00_000 // 創建成功
	NoContent = 204_00_00  // 成功不返回資訊
	//InvalidArgument

	// 400 status code from here
	BadRequest = 400_00_00 // 請求錯誤
	// AccountPasswordError    = 400_0001 // 密碼錯誤
	NotAllowChangeOrderTime        = 400_02_001 //不允許修改訂單時間
	NotAllowSmallerOriginOrderTime = 400_02_002 // 修改時間不可小於原訂單時間

	// 支付相關錯誤碼
	PaymentUnprocessed             = 400_02_101 // 訂單款項未付款
	PaymentProcessing              = 400_02_102 // 訂單款項處理中
	PaymentFailed                  = 400_02_103 // 訂單款項付款失敗
	NoUnPayOrderLineItem           = 400_02_104 //無可付款訂單
	OrderStatusInvalid             = 400_02_201 // 訂單狀態限制不可操作
	NoUnpaidOrderLineItemOrDeposit = 400_02_300 // 無須付款費用項目或押金項目

	// 合約
	ConsignContractEndedAtSmallerStartedAt = 400_04_001 // 合約結束時間小於開始時間
	InvalidConsignContractStartedAt        = 400_04_002 //新合約開始時間不等於就有合約結束時間+1
	ConsignContractIsChanged               = 400_04_003 //合約已被更新
	NonCurrentlyAvailableConsignContract   = 400_04_004 // 非當前有效合約
	ConsignContractIsSuspended             = 400_04_005 //合約已停用
	HasNextContractLostContractDate        = 400_04_006 //有續租合約卻沒有帶續租合約的開始或結束時間
	NextContractStartTimeNotEqualRule      = 400_04_007 // 續租合約的開始時間不符合規則
	AlreadyRenewContract                   = 400_04_008 //已有續約合約
	ContractEndConflict                    = 400_04_009 // 合約結束時間還有訂單沒結束

	// Location
	LocationNoPostalCode = 400_05_001 // 無法取得位置的郵遞區號
	LocationIsDisabled   = 400_05_002 // 該位置已停用

	// 優惠卷
	PromotionSerialAlreadyFetched                     = 400_06_02 // 推廣序號已被使用
	PromotionSerialFetchAllDone                       = 400_06_03 // 推廣序號已被領取完
	PromotionSerialExpired                            = 400_06_04 // 推廣序號已過期
	PromotionSerialAccountNotMatchForPlatformFirstUse = 400_06_05 // 帳號資格不符合平台首用限制

	// 車輛
	VehicleModelClosed                      = 400_07_01  //車款已關閉
	VehicleHadOrder                         = 400_07_02  // 車輛已被訂單
	VehicleClosed                           = 400_07_03  // 車輛已停用
	VehicleNotOnShelf                       = 400_07_04  // 車輛未上架
	VehicleHighlightedNotFound              = 400_07_05  // 車輛精選不存在
	VehicleLocationsCopyConflictOnPickupFee = 400_07_010 // 套用車輛站點設定時，取還車服務費衝突

	// 訂單
	VehicleOrderTimeOverlap                          = 400_10_002 // 車輛訂單時間重疊
	VehicleOrderDownTimeOverlap                      = 400_10_003 // 車輛訂單自用時間重疊
	ChangeOrderDeliveryAndReturnNotInConsignContract = 400_10_004 // 異動車輛時間不符合車輛合約
	OrderHasDepositNotPayed                          = 400_10_005 // 訂單有訂金未付款

	// 證件
	IdentificationIsUnrecognizedByAI                               = 400_11_001 // AI拒絕該證件
	IdentificationHasInvalidBirthDate                              = 400_11_002 // 證件出生日期格式錯誤
	IdentificationHasInvalidIssuedDate                             = 400_11_003 // 證件發行日期格式錯誤
	IdentificationHasInvalidExpiredDate                            = 400_11_004 // 證件有效日期格式錯誤
	IdentificationHasInvalidIDNumber                               = 400_11_005 // 證件編號格式錯誤
	IdentificationUploadIsNotAllowedWhileUnfulfilledOrderUndergoes = 400_11_006 // 有未完成訂單時，不可上傳證件
	MustUploadDrivingLicenseForResidentCertificateUser             = 400_11_010 // 居留證用戶需上傳駕照

	//ETC
	ETCLoginFailure            = 400_12_001 // 遠通ETC登入失敗
	ETCIncorrectLicensePlateNo = 400_12_002 // 遠通ETC查無車牌
	ETCOrdered                 = 400_12_003 // 遠通ETC車輛已有其他訂單
	ETCOrderExist              = 400_12_004 // 遠通ETC已有租車紀錄
	ETCOrderedDuplicateOrderNo = 400_12_005 // 遠通ETC車輛重複訂單編號
	ETCOrderedOverlapPeriod    = 400_12_006 // 遠通ETC車輛重疊租賃時間

	// Rating
	RatingStatusTransitionInvalid = 400_13_001 // 評價狀態轉換不合法
	NegativeRatingMustTagOrText   = 400_13_002 // 負面評價至少需要Tag或文字
	OverRattingTextLimit          = 400_13_003 // 超過評價文字長度限制
	RepairCannotRating            = 400_13_004 // 有車損不能評價
	// 401 status code from here
	Unauthorized = 401_00_00 // 未授權
	// UnusualLogin       = 401_00_01 // 登入異常

	// 403 status code from here
	Forbidden = 403_00_00 // 禁止訪問
	// NoPermission = 403_00_01 // 沒有權限

	// 404 status code from here
	NotFound            = 404_00_00  // 沒有Response
	ConsignTermNotFound = 404_08_001 // 查不到分潤合約
	NoTransactionTask   = 404_14_001 // 沒有待處理的付款工作
	// ResourceNotFound = 404_0001 // 找不到資源

	// 409 status code from here
	Conflict               = 409_00_00  // 衝突
	PaymentProcessConflict = 409_14_001 // 該訂單的付款工作已被其他進程捕獲

	// Content
	ContentVersionConflict = 409_09_001 // 內容版本衝突（提交的 version 與現況不符，需重新載入）

	// 422 status code from here
	UnprocessableEntity        = 422_00_000 // 參數錯誤
	CustomerStatusNoteRequired = 422_01_001 // 客戶狀態備註必填

	// 429 status code from here
	TooManyRequests = 429_00_00 // 請求過多

	// 500 status code from here
	InternalServerError = 500_00_00 // 内部錯誤
	// InvalidPermission   = 500_00_01 // 無效的權限
	UnrecognizableCaptcha = 500_11_001 // 無法辨識驗證碼

	// 501 status code from here
	NotImplemented = 501_00_00 // 功能未實現

	// 503
	ServiceUnavailable     = 503_00_000
	ServiceInMaintenance   = 503_00_001 // 服務維護中
	ServiceForciblyUpdated = 503_00_002 // 強迫更新軟體版本，例如APP
	ServiceSoftlyUpdated   = 503_00_003 // 推薦更新軟體版本，例如APP
)

// HttpCode returns the standard HTTP status code.
func (c CusCode) HttpCode() int {
	var httpCode int

	// Get the http code from the CusCode
	// temporary fix: since there is a legacy typo, code should be 8 digits long but 7 digits version was already released
	if 100_00000 < c && c < 999_99999 {
		httpCode = int(c) / 1_00000
	} else {
		httpCode = int(c) / 1_0000
	}

	// Check if the http code is valid
	if http.StatusText(httpCode) == "" {
		return http.StatusInternalServerError
	}

	return httpCode
}

// GrpcCode returns the corresponding gRPC status code.
func (c CusCode) GrpcCode() codes.Code {
	switch c.HttpCode() {
	case http.StatusOK:
		return codes.OK
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusNotImplemented:
		return codes.Unimplemented
	default:
		return codes.Unknown
	}
}

// Int returns the integer value of the CusCode.
func (c CusCode) Int() int {
	return int(c)
}
