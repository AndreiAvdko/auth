package converter

import (
	"github.com/AndreiAvdko/internal/model"
	modelRepo "github.com/AndreiAvdko/internal/repository/user/model"
)

func ToNoteFromRepo(note *modelRepo.User) *model.Note {
	return &model.Note{
		ID:        note.ID,
		Info:      ToNoteInfoFromRepo(note.Info),
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

func ToNoteInfoFromRepo(info modelRepo.NoteInfo) model.NoteInfo {
	return model.NoteInfo{
		Title:   info.Title,
		Content: info.Content,
	}
}
