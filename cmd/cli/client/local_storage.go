package client

import (
	"encoding/json"
	"gophkeeper/internal/service"
	"io"
	"os"
	"sync"
)

// Storage is an interface of all storage interactions needed for Gophkeeper
type Storage interface {
	UpdatePath(path string) error
	StoreLogoPasses(listLogoPasses []service.LogoPass) ([]service.LogoPass, error)
	StoreTexts(listTexts []service.TextData) ([]service.TextData, error)
	StoreCreditCards(listCreditCards []service.CreditCard) ([]service.CreditCard, error)
	StoreBinaries(binaryList []service.BinaryData) error
	UpdateLogoPass(logoPass service.LogoPass) error
	UpdateText(Text service.TextData) error
	UpdateCreditCard(CreditCard service.CreditCard) error
	UpdateBinaryList(binary service.BinaryData) error
	GetLogoPasses() ([]service.LogoPass, error)
	GetTexts() ([]service.TextData, error)
	GetCreditCards() ([]service.CreditCard, error)
	GetBinaryList() ([]service.BinaryData, error)
}

// FileStorage holds a path for local storage to save files
type FileStorage struct {
	outputPath string
}

// NewStorage creates an instance of a FileStorage along with creating folders for it
func NewStorage(path string) (*FileStorage, error) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return &FileStorage{outputPath: path}, nil
}

// UpdatePath adds a _userLogin postfix to a storage path
func (storage *FileStorage) UpdatePath(path string) error {
	storage.outputPath = path
	err := os.MkdirAll(storage.outputPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// StoreLogoPasses accepts []service.LogoPass acquired from another storage and sorts out which of logo-pass pairs
// have to be updated locally and which of them are newer in this storage. It gathers the ones that are newer here
// and sends []service.LogoPass list for further updating of the remote.
func (storage *FileStorage) StoreLogoPasses(serverLogoPasses []service.LogoPass) ([]service.LogoPass, error) {
	var updLogoPasses []service.LogoPass

	file, err := os.OpenFile(storage.outputPath+logoPassFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return updLogoPasses, err
	}
	defer file.Close()

	var mutex sync.Mutex
	mutex.Lock()

	data, err := io.ReadAll(file)
	if err != nil {
		return updLogoPasses, err
	}

	var storedLogoPasses []service.LogoPass
	if len(data) != 0 {
		err = json.Unmarshal(data, &storedLogoPasses)
		if err != nil {
			return updLogoPasses, err
		}
	}
	for i, serverLogoPass := range serverLogoPasses {
		newEntry := true
		for _, storedLogoPass := range storedLogoPasses {
			if serverLogoPass.Description == storedLogoPass.Description {
				if serverLogoPass.UpdatedAt.After(storedLogoPass.UpdatedAt) {
					storedLogoPasses[i] = serverLogoPass
				} else if serverLogoPass.UpdatedAt.Before(storedLogoPass.UpdatedAt) {
					storedLogoPass.Overwrite = true
					updLogoPasses = append(updLogoPasses, storedLogoPass)
				}
				newEntry = false
				break
			}
		}
		if newEntry {
			storedLogoPasses = append(storedLogoPasses, serverLogoPass)
		}
	}
	for _, storedLogoPass := range storedLogoPasses {
		newEntry := true
		for _, serverLogoPass := range serverLogoPasses {
			if serverLogoPass.Description == storedLogoPass.Description {
				newEntry = false
				break
			}
		}
		if newEntry {
			updLogoPasses = append(updLogoPasses, storedLogoPass)
		}
	}

	if len(storedLogoPasses) != 0 {
		jsonBytes, err := json.Marshal(storedLogoPasses)
		if err != nil {
			return updLogoPasses, err
		}
		err = os.WriteFile(storage.outputPath+logoPassFile, jsonBytes, 0644)
		if err != nil {
			return updLogoPasses, err
		}
	}
	mutex.Unlock()

	return updLogoPasses, nil
}

// StoreTexts accepts []service.TextData acquired from another storage and sorts out which of text data entries
// have to be updated locally and which of them are newer in this storage. It gathers the ones that are newer here
// and sends []service.TextData list for further updating of the remote.
func (storage *FileStorage) StoreTexts(serverTexts []service.TextData) ([]service.TextData, error) {
	var updTexts []service.TextData

	file, err := os.OpenFile(storage.outputPath+textFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return updTexts, err
	}
	defer file.Close()

	var mutex sync.Mutex
	mutex.Lock()

	data, err := io.ReadAll(file)
	if err != nil {
		return updTexts, err
	}

	var storedTexts []service.TextData
	if len(data) != 0 {
		err = json.Unmarshal(data, &storedTexts)
		if err != nil {
			return updTexts, err
		}
	}
	for _, serverText := range serverTexts {
		newEntry := true
		for i, storedText := range storedTexts {
			if serverText.Description == storedText.Description {
				if serverText.UpdatedAt.After(storedText.UpdatedAt) {
					storedTexts[i] = serverText
				} else if serverText.UpdatedAt.Before(storedText.UpdatedAt) {
					storedText.Overwrite = true
					updTexts = append(updTexts, storedText)
				}
				newEntry = false
				break
			}
		}
		if newEntry {
			storedTexts = append(storedTexts, serverText)
		}
	}
	for _, storedText := range storedTexts {
		newEntry := true
		for _, serverText := range serverTexts {
			if serverText.Description == storedText.Description {
				newEntry = false
				break
			}
		}
		if newEntry {
			updTexts = append(updTexts, storedText)
		}
	}

	if len(storedTexts) != 0 {
		jsonBytes, err := json.Marshal(storedTexts)
		if err != nil {
			return updTexts, err
		}
		err = os.WriteFile(storage.outputPath+textFile, jsonBytes, 0644)
		if err != nil {
			return updTexts, err
		}
	}
	mutex.Unlock()

	return updTexts, nil
}

// StoreCreditCards accepts []service.CreditCard acquired from another storage and sorts out which of credit cards entries
// have to be updated locally and which of them are newer in this storage. It gathers the ones that are newer here
// and sends []service.CreditCard list for further updating of the remote.
func (storage *FileStorage) StoreCreditCards(serverCreditCards []service.CreditCard) ([]service.CreditCard, error) {
	var updCreditCards []service.CreditCard

	file, err := os.OpenFile(storage.outputPath+creditCardFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return updCreditCards, err
	}
	defer file.Close()

	var mutex sync.Mutex
	mutex.Lock()

	data, err := io.ReadAll(file)
	if err != nil {
		return updCreditCards, err
	}

	var storedCreditCards []service.CreditCard
	if len(data) != 0 {
		err = json.Unmarshal(data, &storedCreditCards)
		if err != nil {
			return updCreditCards, err
		}
	}
	for i, serverCard := range serverCreditCards {
		newEntry := true
		for _, storedCard := range storedCreditCards {
			if serverCard.Number == storedCard.Number {
				if serverCard.UpdatedAt.After(storedCard.UpdatedAt) {
					storedCreditCards[i] = serverCard
				} else if serverCard.UpdatedAt.Before(storedCard.UpdatedAt) {
					storedCard.Overwrite = true
					updCreditCards = append(updCreditCards, storedCard)
				}
				newEntry = false
				break
			}
		}
		if newEntry {
			storedCreditCards = append(storedCreditCards, serverCard)
		}
	}
	for _, storedCard := range storedCreditCards {
		newEntry := true
		for _, serverCard := range serverCreditCards {
			if serverCard.Number == storedCard.Number {
				newEntry = false
				break
			}
		}
		if newEntry {
			updCreditCards = append(updCreditCards, storedCard)
		}
	}

	if len(storedCreditCards) != 0 {
		jsonBytes, err := json.Marshal(storedCreditCards)
		if err != nil {
			return updCreditCards, err
		}
		err = os.WriteFile(storage.outputPath+creditCardFile, jsonBytes, 0644)
		if err != nil {
			return updCreditCards, err
		}
	}
	mutex.Unlock()

	return updCreditCards, nil

}

// StoreBinaries accepts []service.BinaryData acquired from another storage and sorts out which of text data entries
// have to be updated locally and which of them are newer in this storage. It gathers the ones that are newer here
// and sends []service.BinaryData list for further updating of the remote.
//
// **NOTE** that it only stores the descriptions of binaries, not binaries themselves. Those are downloaded by request.
func (storage *FileStorage) StoreBinaries(serverBinaries []service.BinaryData) error {
	file, err := os.OpenFile(storage.outputPath+binaryListFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var mutex sync.Mutex
	mutex.Lock()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var storedBinaries []service.BinaryData
	if len(data) != 0 {
		err = json.Unmarshal(data, &storedBinaries)
		if err != nil {
			return err
		}
	}
	for i, serverBinary := range serverBinaries {
		newEntry := true
		for _, storedBinary := range storedBinaries {
			if serverBinary.Description == storedBinary.Description {
				newEntry = false
				if serverBinary.UpdatedAt.After(storedBinary.UpdatedAt) {
					storedBinaries[i] = serverBinary
					break
				}
			}
		}
		if newEntry {
			storedBinaries = append(storedBinaries, serverBinary)
		}
	}

	if len(storedBinaries) != 0 {
		jsonBytes, err := json.Marshal(storedBinaries)
		if err != nil {
			return err
		}
		err = os.WriteFile(storage.outputPath+binaryListFile, jsonBytes, 0644)
		if err != nil {
			return err
		}
	}
	mutex.Unlock()

	return nil
}

// UpdateLogoPass accepts service.LogoPass and writes it to local storage if it's new.
// If it's not, checks the 'overwrite' flag to figure out if the data should be updated or not.
func (storage *FileStorage) UpdateLogoPass(logoPass service.LogoPass) error {
	var mutex sync.Mutex
	mutex.Lock()

	file, err := os.OpenFile(storage.outputPath+logoPassFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var listLogoPasses []service.LogoPass
	if len(data) != 0 {
		err = json.Unmarshal(data, &listLogoPasses)
		if err != nil {
			return err
		}

		newEntry := true
		var newID uint = 1
		for i, existingLogoPass := range listLogoPasses {
			if existingLogoPass.Description == logoPass.Description {
				newEntry = false
				if logoPass.Overwrite {
					listLogoPasses[i] = logoPass
					break
				} else {
					return ErrAlreadyExists
				}
			}
			newID = existingLogoPass.ID + 1
		}
		if newEntry {
			logoPass.ID = newID
			listLogoPasses = append(listLogoPasses, logoPass)
		}
	} else {
		logoPass.ID = 1
		listLogoPasses = append(listLogoPasses, logoPass)
	}

	jsonBytes, err := json.Marshal(listLogoPasses)
	if err != nil {
		return err
	}
	err = os.WriteFile(storage.outputPath+logoPassFile, jsonBytes, 0644)
	if err != nil {
		return err
	}
	mutex.Unlock()

	return nil
}

// UpdateText accepts service.TextData and writes it to local storage if it's new.
// If it's not, checks the 'overwrite' flag to figure out if the data should be updated or not.
func (storage *FileStorage) UpdateText(text service.TextData) error {
	var mutex sync.Mutex
	mutex.Lock()
	file, err := os.OpenFile(storage.outputPath+textFile, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var listTexts []service.TextData
	if len(data) != 0 {
		err = json.Unmarshal(data, &listTexts)
		if err != nil {
			return err
		}

		newEntry := true
		var newID uint = 1
		for i, existingText := range listTexts {
			if existingText.Description == text.Description {
				newEntry = false
				if text.Overwrite {
					listTexts[i] = text
					break
				} else {
					return ErrAlreadyExists
				}
			}
			newID = existingText.ID + 1
		}
		if newEntry {
			text.ID = newID
			listTexts = append(listTexts, text)
		}
	} else {
		text.ID = 1
		listTexts = append(listTexts, text)
	}

	jsonBytes, err := json.Marshal(listTexts)
	if err != nil {
		return err
	}
	err = os.WriteFile(storage.outputPath+textFile, jsonBytes, 0644)
	if err != nil {
		return err
	}
	mutex.Unlock()

	return nil
}

// UpdateCreditCard accepts service.CreditCard and writes it to local storage if it's new.
// If it's not, checks the 'overwrite' flag to figure out if the data should be updated or not.
func (storage *FileStorage) UpdateCreditCard(creditCard service.CreditCard) error {
	var mutex sync.Mutex
	mutex.Lock()

	file, err := os.OpenFile(storage.outputPath+creditCardFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var listCreditCards []service.CreditCard
	if len(data) != 0 {
		err = json.Unmarshal(data, &listCreditCards)
		if err != nil {
			return err
		}

		newEntry := true
		var newID uint = 1
		for i, existingCard := range listCreditCards {
			if existingCard.Number == creditCard.Number {
				newEntry = false
				if creditCard.Overwrite {
					listCreditCards[i] = creditCard
					break
				} else {
					return ErrAlreadyExists
				}
			}
			newID = existingCard.ID + 1
		}
		if newEntry {
			creditCard.ID = newID
			listCreditCards = append(listCreditCards, creditCard)
		}
	} else {
		creditCard.ID = 1
		listCreditCards = append(listCreditCards, creditCard)
	}

	jsonBytes, err := json.Marshal(listCreditCards)
	if err != nil {
		return err
	}
	err = os.WriteFile(storage.outputPath+creditCardFile, jsonBytes, 0644)
	if err != nil {
		return err
	}
	mutex.Unlock()

	return nil
}

// UpdateBinaryList accepts service.BinaryData and writes it to local storage if it's new.
// If it's not, checks the 'overwrite' flag to figure out if the data should be updated or not.
func (storage *FileStorage) UpdateBinaryList(binary service.BinaryData) error {
	var mutex sync.Mutex
	mutex.Lock()

	file, err := os.OpenFile(storage.outputPath+binaryListFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var binaryList []service.BinaryData
	if len(data) != 0 {
		err = json.Unmarshal(data, &binaryList)
		if err != nil {
			return err
		}

		newEntry := true
		var newID uint = 1
		for i, existingBinary := range binaryList {
			if existingBinary.Description == binary.Description {
				newEntry = false
				if binary.Overwrite {
					binaryList[i] = binary
					break
				} else {
					return ErrAlreadyExists
				}
			}
			newID = existingBinary.ID + 1
		}
		if newEntry {
			binary.ID = newID
			binaryList = append(binaryList, binary)
		}
	} else {
		binary.ID = 1
		binaryList = append(binaryList, binary)
	}

	jsonBytes, err := json.Marshal(binaryList)
	if err != nil {
		return err
	}
	err = os.WriteFile(storage.outputPath+binaryListFile, jsonBytes, 0644)
	if err != nil {
		return err
	}
	mutex.Unlock()

	return nil
}

// GetLogoPasses returns the list of all service.LogoPass entries available locally
func (storage *FileStorage) GetLogoPasses() ([]service.LogoPass, error) {
	var listLogoPasses []service.LogoPass

	file, err := os.OpenFile(storage.outputPath+logoPassFile, os.O_RDONLY, 0644)
	if err != nil {
		return listLogoPasses, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return listLogoPasses, err
	}
	if len(data) != 0 {
		err = json.Unmarshal(data, &listLogoPasses)
		if err != nil {
			return listLogoPasses, err
		}
	}
	return listLogoPasses, nil
}
func (storage *FileStorage) GetTexts() ([]service.TextData, error) {
	var listTexts []service.TextData

	file, err := os.OpenFile(storage.outputPath+textFile, os.O_RDONLY, 0644)
	if err != nil {
		return listTexts, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return listTexts, err
	}
	if len(data) != 0 {
		err = json.Unmarshal(data, &listTexts)
		if err != nil {
			return listTexts, err
		}
	}
	return listTexts, nil
}

// GetCreditCards returns the list of all service.CreditCard entries available locally
func (storage *FileStorage) GetCreditCards() ([]service.CreditCard, error) {
	var listCreditCards []service.CreditCard

	file, err := os.OpenFile(storage.outputPath+creditCardFile, os.O_RDONLY, 0644)
	if err != nil {
		return listCreditCards, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return listCreditCards, err
	}

	if len(data) != 0 {
		err = json.Unmarshal(data, &listCreditCards)
		if err != nil {
			return listCreditCards, err
		}
	}
	return listCreditCards, nil
}

// GetBinaryList returns the list of all service.BinaryData entries available locally
func (storage *FileStorage) GetBinaryList() ([]service.BinaryData, error) {
	var BinaryList []service.BinaryData

	file, err := os.OpenFile(storage.outputPath+binaryListFile, os.O_RDONLY, 0644)
	if err != nil {
		return BinaryList, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return BinaryList, err
	}
	if len(data) != 0 {
		err = json.Unmarshal(data, &BinaryList)
		if err != nil {
			return BinaryList, err
		}
	}
	return BinaryList, nil
}

// ClearAll destroys local storage
func (storage *FileStorage) ClearAll() error {
	err := os.RemoveAll(storage.outputPath)
	if err != nil {
		return err
	}
	return nil
}
