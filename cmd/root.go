package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"adonis/core"
)

var filesPath string
var rootCmd = &cobra.Command{
	Use:   "adonis",
	Short: "Analyze resource list, and process the image inside",
	Long: `< Analyze resource list, and process the image inside >

   /::\  \       /::\  \       /::\  \        \:\  \       ___         /:/ _/_   
  /:/\:\  \     /:/\:\  \     /:/\:\  \        \:\  \     /\__\       /:/ /\  \  
 /:/ /::\  \   /:/  \:\__\   /:/  \:\  \   _____\:\  \   /:/__/      /:/ /::\  \ 
/:/_/:/\:\__\ /:/__/ \:|__| /:/__/ \:\__\ /::::::::\__\ /::\  \     /:/_/:/\:\__\
\:\/:/  \/__/ \:\  \ /:/  / \:\  \ /:/  / \:\~~\~~\/__/ \/\:\  \__  \:\/:/ /:/  /
 \::/__/       \:\  /:/  /   \:\  /:/  /   \:\  \        ~~\:\/\__\  \::/ /:/  / 
  \:\  \        \:\/:/  /     \:\/:/  /     \:\  \          \::/  /   \/_/:/  /  
   \:\__\        \::/  /       \::/  /       \:\__\         /:/  /      /:/  /   
	`,
}

var parseCmd = &cobra.Command{
	Use:     "parse",
	Short:   "Find all images in yaml where the key is 'image'",
	Example: "  adonis parse -f <file_path/dir_path>",
	Run: func(cmd *cobra.Command, args []string) {
		core.GetImages(filesPath)
	},
}

var pullCmd = &cobra.Command{
	Use:     "pull",
	Short:   "Pull all images in yaml where the key is 'image'",
	Example: "  adonis pull -f <file_path/dir_path>",
	Run: func(cmd *cobra.Command, args []string) {
		imges := core.GetImages(filesPath)
		fmt.Printf("\n")

		core.PullImages(core.CreateDockerClient(), imges)
	},
}

var savePath string
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save all images in yaml where the key is 'image'",
	Example: `  adonis save -f <file_path/dir_path>
  adonis save -f <file_path/dir_path> -p <image_save_path>`,
	Run: func(cmd *cobra.Command, args []string) {
		imges := core.GetImages(filesPath)
		fmt.Printf("\n")

		core.PullImages(core.CreateDockerClient(), imges)
		fmt.Printf("\n")

		core.SaveImages(core.CreateDockerClient(), imges, savePath)
	},
}

var repoPath string
var isDeleteOriginTag bool
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Pull all images in yaml where the key is 'image', and tag them",
	Example: `  adonis tag -f <file_path/dir_path> -r <new_repository_path>
  adonis tag -f <file_path/dir_path> -r <new_repository_path> -d`,
	Run: func(cmd *cobra.Command, args []string) {
		cli := core.CreateDockerClient()
		images := core.GetImages(filesPath)
		fmt.Printf("\n")

		core.PullImages(cli, images)
		fmt.Printf("\n")

		core.TagImages(cli, images, repoPath)
		fmt.Printf("\n")

		if isDeleteOriginTag {
			core.RemoveImages(cli, images)
			fmt.Printf("\n")
		}
	},
}

var isSaveImage bool
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push all images in yaml where the key is 'image', tag and push them",
	Example: `  adonis push -f <file_path/dir_path> -r <new_repository_path>
  adonis push -f <file_path/dir_path> -r <new_repository_path> -d
  adonis push -f <file_path/dir_path> -r <new_repository_path> -s
  adonis push -f <file_path/dir_path> -r <new_repository_path> -s -p <image_save_path>
  adonis push -f <file_path/dir_path> -r <new_repository_path> -s -p <image_save_path> -d`,
	Run: func(cmd *cobra.Command, args []string) {
		cli := core.CreateDockerClient()
		images := core.GetImages(filesPath)
		fmt.Printf("\n")

		core.PullImages(cli, images)
		fmt.Printf("\n")

		if isSaveImage {
			core.SaveImages(cli, images, savePath)
			fmt.Printf("\n")
		}

		if isDeleteOriginTag {
			core.RemoveImages(cli, images)
			fmt.Printf("\n")
		}

		newImage := core.TagImages(cli, images, repoPath)
		fmt.Printf("\n")

		core.PushImages(cli, newImage)
	},
}

func initFlag() {
	rootCmd.PersistentFlags().StringVarP(&filesPath, "files", "f", "", "yaml to be parsed")
	rootCmd.MarkPersistentFlagRequired("files")

	rootCmd.AddCommand(parseCmd)

	rootCmd.AddCommand(pullCmd)

	saveCmd.Flags().StringVarP(&savePath, "path", "p", "./", "path to save image")
	rootCmd.AddCommand(saveCmd)

	tagCmd.Flags().StringVarP(&repoPath, "repo", "r", "", "path to new repository")
	tagCmd.Flags().BoolVarP(&isDeleteOriginTag, "delete", "d", false, "delete origin tag or not")
	tagCmd.MarkFlagRequired("repo")
	rootCmd.AddCommand(tagCmd)

	pushCmd.Flags().StringVarP(&savePath, "path", "p", "./", "path to save image")
	pushCmd.Flags().StringVarP(&repoPath, "repo", "r", "", "path to new repository")
	pushCmd.Flags().BoolVarP(&isSaveImage, "save", "s", false, "save image or not")
	pushCmd.Flags().BoolVarP(&isDeleteOriginTag, "delete", "d", false, "delete origin tag or not")
	pushCmd.MarkFlagRequired("repo")
	rootCmd.AddCommand(pushCmd)
}

func Execute() {
	initFlag()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
