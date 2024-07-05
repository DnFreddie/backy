package trip

import (
	"fmt"
	"os"
	"sync"
	"github.com/spf13/cobra"
	"github.com/DnFreddie/backy/utils"
	"gorm.io/gorm"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Indexes the directory and adds it to the existing index group for future scanning",
	Long: `
	If there's no path specified for scanning, it adds the directory to the pool.
	It then checks each file, computes checksums, and stores them in the database for comparison during subsequent scans.
`,

	Run: func(cmd *cobra.Command, args []string) {
		if addPath != "" {
			_, err := os.Stat(addPath)
			if err != nil {
				fmt.Println("the driectry doesn't exist")
				os.Exit(1)

			}
			tripAdd(addPath)
		}

	},
}

var (
	addPath string
)

func init() {
	addCmd.Flags().StringVarP(&addPath, "path", "p", "", "the path to add to the scan")
	addCmd.MarkFlagRequired("path")

}




func tripAdd(fPath string) error {
	db, err := utils.InitDb(DB_PATH, &utils.FileProps{})
	if err != nil {
		return err
	}

	isNew, err := CreateConfig(fPath)
	if err != nil {
		return err
	}

	if isNew {
		fmt.Printf("The %v does already exist in db. Try scan flag\n", fPath)
		return nil
	}

	var wg sync.WaitGroup
	ch := make(chan utils.FileProps)
	numWorkers := 1

	wg.Add(1)
	go func() {
		defer wg.Done()
		ScanRecursivly(fPath, db, ch)
		close(ch)
	}()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processBatches(ch, db)
		}()
	}

	wg.Wait()
	return nil
}

func processBatches(ch chan utils.FileProps, db *gorm.DB) {
	var batch []utils.FileProps
	var count int
	for item := range ch {
		count++
		batch = append(batch, item)
		if len(batch) == 100 {
			db.CreateInBatches(batch, len(batch))
			fmt.Println("Processed count during execution:", count)
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		db.CreateInBatches(batch, len(batch))
		fmt.Println("Processed count during execution:", count)
	}
}
