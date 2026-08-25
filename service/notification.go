package service

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"labops/model"
	"labops/storage"
)

type Notice struct {
	ID        string
	RecordID  string
	Recipient string
	Subject   string
	Body      string
	CreatedAt time.Time
	Read      bool
}

type Notifier struct {
	Store  *storage.Store
	Clock  func() time.Time
	mu     sync.Mutex
	notice []Notice
}

func NewNotifier(store *storage.Store) *Notifier {
	return &Notifier{Store: store, Clock: time.Now, notice: make([]Notice, 0)}
}

func (n *Notifier) Notify(record model.Record, recipient, reason string) Notice {
	n.mu.Lock()
	defer n.mu.Unlock()
	notice := Notice{
		ID:        fmt.Sprintf("notice-%s-%d", record.ID, len(n.notice)+1),
		RecordID:  record.ID,
		Recipient: recipient,
		Subject:   fmt.Sprintf("设备记录 %s 状态更新", record.ID),
		Body:      fmt.Sprintf("%s: %s", reason, model.StatusLabel(record.Status)),
		CreatedAt: n.Clock(),
	}
	n.notice = append(n.notice, notice)
	return notice
}

func (n *Notifier) MarkRead(id string) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	for index := range n.notice {
		if n.notice[index].ID == id {
			n.notice[index].Read = true
			return true
		}
	}
	return false
}

func (n *Notifier) List(recipient string, unreadOnly bool) []Notice {
	n.mu.Lock()
	defer n.mu.Unlock()
	items := make([]Notice, 0)
	for _, item := range n.notice {
		if recipient != "" && item.Recipient != recipient {
			continue
		}
		if unreadOnly && item.Read {
			continue
		}
		items = append(items, item)
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	return items
}

func (n *Notifier) UnreadCount(recipient string) int {
	return len(n.List(recipient, true))
}

func (n *Notifier) Clear() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	count := len(n.notice)
	n.notice = make([]Notice, 0)
	return count
}

func (n *Notifier) Has(id string) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, item := range n.notice {
		if item.ID == id {
			return true
		}
	}
	return false
}

func (n *Notifier) Latest(recordID string) (Notice, bool) {
	items := n.List("", false)
	for _, item := range items {
		if item.RecordID == recordID {
			return item, true
		}
	}
	return Notice{}, false
}

func (n *Notifier) NotifyTransition(record model.Record, actor string) Notice {
	reason := fmt.Sprintf("%s changed progress by %s", actor, actor)
	return n.Notify(record, record.Assignee, reason)
}

func (n *Notifier) ValidateNotice(notice Notice) error {
	if notice.ID == "" || notice.RecordID == "" {
		return fmt.Errorf("notice identity required")
	}
	if notice.Recipient == "" {
		return fmt.Errorf("recipient required")
	}
	if notice.Subject == "" || notice.Body == "" {
		return fmt.Errorf("notice content required")
	}
	return nil
}
