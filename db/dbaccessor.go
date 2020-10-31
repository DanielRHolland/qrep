package db

import qrep "github.com/DanielRHolland/qrep/models"

type DbAccessor interface {
    InsertItem(qrep.TrackedItemType) (string, error)

    GetItem(string) (qrep.TrackedItemType, error)

    AddIssueToItem(qrep.IssueType, string) error

    GetTrackedItems(int) ([]qrep.TrackedItemType, error)

    SearchTrackedItems(int,string) ([]qrep.TrackedItemType, error)

    UpdateDbIssue(qrep.IssueType) error

    GetItemsFromIds([]string) ([]qrep.TrackedItemType, error)

    RemoveItemsFromDb([]string) error


}
