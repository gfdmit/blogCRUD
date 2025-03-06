package app

import (
	"blog/config"
	"blog/internal/model"
	"blog/internal/repository"
	"fmt"
)

func Run(conf config.Config) error {
	repo, err := repository.New(conf)
	if err != nil {
		return fmt.Errorf("repository.New: %v", err)
	}
	if err := repo.Migrate(&model.Post{}); err != nil {
		return fmt.Errorf("repo.Migrate: %v", err)
	}
	for {
		var opt int
		fmt.Print("Menu:\n",
			"  1. Create\n",
			"  2. Read\n",
			"  3. Update\n",
			"  4. Delete\n",
			"  9. Exit\n",
			"Choose action: ")

		if _, err := fmt.Scanf("%d\n", &opt); err != nil {
			var temp string
			fmt.Scanln(&temp)
		}

		switch opt {
		case 1:
			repo.Create()
		case 2:
			repo.Read()
		case 3:
			repo.Update()
		case 4:
			repo.Delete()
		case 9:
			return nil
		default:
			fmt.Println("Incorrect input")
		}
	}
}
