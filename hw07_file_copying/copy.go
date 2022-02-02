package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOffsetValue           = errors.New("negative offset")
	ErrLimitValue            = errors.New("negative limit")
	// MaxCopyPerStep максимальное кол-во байт за итерацию.
	MaxCopyPerStep int64 = 128
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		return ErrOffsetValue
	}

	if limit < 0 {
		return ErrLimitValue
	}

	// открываем файл на чтение
	fileSrc, err := os.OpenFile(fromPath, os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	// не забываем закрывать
	defer fileSrc.Close()

	// открываем файл на запись
	fileDst, err := os.Create(toPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	// не забываем закрывать
	defer fileDst.Close()

	// проверяем длину файла
	statSrc, err := fileSrc.Stat()
	if err != nil {
		return err
	}
	// если offset больше размера файла, выводим ошибку
	if offset > statSrc.Size() {
		return ErrOffsetExceedsFileSize
	}

	// устанавливаем курсор в положение offset от начала файла
	if _, err = fileSrc.Seek(offset, io.SeekStart); err != nil {
		return ErrUnsupportedFile
	}

	if limit == 0 {
		// если лимит не задан, рассчитываем его до конца файла
		limit = statSrc.Size() - offset
	} else {
		// если задан, то уменьшаем до кол-ва байт до конца файла
		limit = min(statSrc.Size()-offset, limit)
	}

	// копируем limit байт

	// прогресс-бар
	stepCount := limit / MaxCopyPerStep
	bar := pb.New(int(stepCount) + 1)
	defer bar.Finish()
	bar.Start()

	var copiedByte int64
	for copiedByte < limit {
		// считаем сколько байт еще нужно прочитать, но не более константы MaxCopyPerStep
		count := min(MaxCopyPerStep, limit-copiedByte)
		// копируем count байт
		n, err := io.CopyN(fileDst, fileSrc, count)
		if err != nil {
			if errors.Is(err, io.EOF) {
				// если конец файла, то выходим
				break
			}
			return ErrUnsupportedFile
		}
		// добавляем кол-во прочитанных байт
		copiedByte += n
		// инкрементируем прогресс-бар
		bar.Increment()
	}

	// пришли к успеху )) всё ок
	return nil
}

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
