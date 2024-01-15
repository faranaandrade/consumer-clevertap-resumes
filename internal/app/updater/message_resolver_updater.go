package updater

import (
	"context"
	"strings"
	"time"

	"github.com/occmundial/go-common/logger"

	"github.com/occmundial/consumer-clevertap-resumes/internal/models"
)

const typeUse = "event"
const evtName = "success_apply"

func NewMessageResolverUpdater(clevetapRepository *APIClevetapRepository, dataUserRepository *DataUserRepositoryDatabase,
	log *logger.Log) *MessageResolverUpdater {
	return &MessageResolverUpdater{
		ClevetapRepository: clevetapRepository,
		DataUserRepository: dataUserRepository,
		log:                log,
	}
}

type MessageResolverUpdater struct {
	PageSize           int
	ClevetapRepository ClevertapGetter
	DataUserRepository DataUserGetter
	log                logger.Logger
}

func (resolver *MessageResolverUpdater) Process(ctx context.Context, inputMsg *models.MessageToProcess) error {
	if err := resolver.ClevetapRepository.APICheck(); err != nil {
		return err
	}
	if strings.Contains(inputMsg.UserID, "xmx") {
		resolver.log.Infof("it's recruiter accounts, it will not be processed")
		return nil
	}
	email, err := resolver.DataUserRepository.GetDBInfo(ctx, inputMsg.UserID)
	if err != nil {
		resolver.log.Error("updater", "Process", err)
		return err
	}
	outputMsg := resolver.mapEmailStatisticForCreation(inputMsg, email)
	err = resolver.ClevetapRepository.SendRequest(outputMsg)
	if err != nil {
		resolver.log.Error("updater", "Process", err)
		return err
	}
	return nil
}

func (resolver *MessageResolverUpdater) mapEmailStatisticForCreation(inputMsg *models.MessageToProcess, email string) *models.ClevetapBody {
	evtData := models.EvtData{
		ResumeID:          inputMsg.ResumeID,
		Email:             email,
		CvReady:           inputMsg.CvReady,
		CvAttached:        inputMsg.CvAttached,
		EducationLevel:    inputMsg.EducationLevel,
		YearsOfExperience: inputMsg.YearsOfExperience,
	}

	data := models.ClevertapData{
		Identity: email,
		TS:       time.Now().Unix(),
		TypeUse:  typeUse,
		EvtName:  evtName,
		EvtData:  evtData,
	}
	var body []models.ClevertapData
	return &models.ClevetapBody{
		Body: append(body, data),
	}
}

func (resolver *MessageResolverUpdater) IsValid(message *models.MessageToProcess) bool {
	return len(message.UserID) > 0
}

func (resolver *MessageResolverUpdater) GetRetryNumber(message *models.MessageToProcess) int {
	return message.RetryNumber
}
func (resolver *MessageResolverUpdater) SetRetryNumber(message *models.MessageToProcess, value int) {
	message.RetryNumber = value
}
