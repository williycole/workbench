package main

type notification interface {
	importance() int
}

type directMessage struct {
	senderUsername string
	messageContent string
	priorityLevel  int
	isUrgent       bool
}

func (dm directMessage) importance() int {
	if dm.isUrgent {
		return 50
	}
	return dm.priorityLevel
}

type groupMessage struct {
	groupName      string
	messageContent string
	priorityLevel  int
}

func (gm groupMessage) importance() int {
	return gm.priorityLevel
}

type systemAlert struct {
	alertCode      string
	messageContent string
}

func (sa systemAlert) importance() int {
	return 100
}

// Complete the processNotification function. It should identify the type and return different values for each type
// If the notification does not match any known type, return an empty string and a score of -1.
// - For a directMessage, return the sender's username and importance score.
// - For a groupMessage, return the group's name and the importance score.
// - For an systemAlert, return the alert code and the importance score.
func processNotification(n notification) (string, int) {
	switch v := n.(type) {
	case *directMessage:
		return v.senderUsername, v.importance()
	case *groupMessage:
		return v.groupName, v.importance()
	case *systemAlert:
		return v.alertCode, v.importance()
	default:
		return "", 0
	}
}
