package database

func (record BookRecord)Equal(compare BookRecord) bool {
	return record.Site == compare.Site && record.Id == compare.Id &&
		record.HashCode == compare.HashCode && record.Title == compare.Title &&
		record.WriterId == compare.WriterId && record.Type == compare.Type &&
		record.UpdateDate == compare.UpdateDate &&
		record.UpdateChapter == compare.UpdateChapter &&
		record.Status == compare.Status
}

func (record WriterRecord)Equal(compare WriterRecord) bool {
	return record.Id == compare.Id && record.Name == compare.Name
}

func (record ErrorRecord)Equal(compare ErrorRecord) bool {
	return record.Site == compare.Site && record.Id == compare.Id &&
		record.Error.Error() == compare.Error.Error()
}