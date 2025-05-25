package snippets

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/maycon-jesus/mj-cli/utils/terminal"
	"github.com/spf13/cobra"
)

var PackagesRpm = []string{
	"gnome-shell-extension-gsconnect",
	"gnome-tweaks",
	"code",
	"zsh",
	"make",
	"zsh",
}

var PackagesRpmDownload = []string{
	"https://dl.google.com/linux/direct/google-chrome-stable_current_x86_64.rpm",
	"https://downloads.1password.com/linux/rpm/stable/x86_64/1password-latest.rpm",
	"https://desktop.docker.com/linux/main/amd64/docker-desktop-x86_64.rpm",
}

var PackagesFlatpak = []string{
	"com.discordapp.Discord",
	"io.missioncenter.MissionCenter",
	"com.spotify.Client",
	"md.obsidian.Obsidian",
	"io.github.flattool.Warehouse",
	"org.gnome.Extensions",
	"com.mattjakeman.ExtensionManager",
	"com.surfshark.Surfshark",
	"org.gnome.Mahjongg",
	"org.telegram.desktop",
	"io.github.fabrialberio.pinapp",
}

var NerdFonts = []string{
	"Hack",
	"FiraCode",
	"Ubuntu",
	"UbuntuMono",
	"UbuntuSans",
	"JetBrainsMono",
}

var AppImageDir = "Apps"

//go:embed gitconfig.txt
var gitConfigData []byte

//go:embed sshconfig.txt
var sshConfigData []byte

//go:embed 1password-config.txt
var onePasswordConfigData []byte

var PostInstallFedoraCommand = &cobra.Command{
	Use: "post-install-fedora",
	Run: RunPostInstallFedoraCommand,
}

func GetPostInstallFedoraCommand() *cobra.Command {
	return PostInstallFedoraCommand
}

var terminalOptions = terminal.RunCommandOptions{
	Debug: true,
}

func printGroupName(groupName string) {
	fmt.Println("====[" + groupName + "]====")
}

func printGroupCommand(cmd string) {
	fmt.Println("==> " + cmd)
}

func execCommandsGroup(groupName string, commands []string) {
	printGroupName(groupName)
	for _, command := range commands {
		printGroupCommand("Executando o comando: " + command)
		err := terminal.RunCommandRealtime(command, terminalOptions)
		if err != nil {
			printGroupCommand("Erro ao executar o comando")
			cobra.CheckErr(err)
			return
		}
	}
}

func RunPostInstallFedoraCommand(cmd *cobra.Command, args []string) {
	tmpDir, _ := os.MkdirTemp(os.TempDir(), "mj-cli-")
	homeDir, _ := os.UserHomeDir()
	//defer func() {
	//	_ = os.RemoveAll(tmpDir)
	//}()

	// Atualizando o sistema
	execCommandsGroup("ATUALIZANDO O SISTEMA", []string{
		"sudo dnf update -y",
		"flatpak update -y",
	})

	// Atualizando o sistema
	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()
		writer.Write([]byte("[code]\nname=Visual Studio Code\nbaseurl=https://packages.microsoft.com/yumrepos/vscode\nenabled=1\nautorefresh=1\ntype=rpm-md\ngpgcheck=1\ngpgkey=https://packages.microsoft.com/keys/microsoft.asc"))
	}()
	terminalOptions.Stdin = reader
	execCommandsGroup("CONFIGURANDO REPOSITÓRIOS -> VsCode", []string{
		"sudo tee /etc/yum.repos.d/vscode.repo",
	})
	terminalOptions.Stdin = os.Stdin

	execCommandsGroup("CONFIGURANDO REPOSITÓRIOS -> Docker", []string{
		"sudo dnf -y install dnf-plugins-core",
		"sudo dnf-3 config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo",
	})

	rpmPackagesFormated := strings.Join(PackagesRpm, " ")
	execCommandsGroup("INSTALANDO PACOTES RPM", []string{
		fmt.Sprintf("sudo dnf install -y %s", rpmPackagesFormated),
	})

	cmdsDownloadAndInstallRpm := []string{}
	for _, packageUrl := range PackagesRpmDownload {
		cmd := fmt.Sprintf("wget %s -O %s/package.rpm", packageUrl, tmpDir)
		cmdsDownloadAndInstallRpm = append(cmdsDownloadAndInstallRpm, cmd)
		cmdsDownloadAndInstallRpm = append(cmdsDownloadAndInstallRpm, "sudo dnf install -y "+tmpDir+"/package.rpm")
		cmdsDownloadAndInstallRpm = append(cmdsDownloadAndInstallRpm, "rm -rf "+tmpDir+"/package.rpm")
	}
	execCommandsGroup("BAIXANDO E INSTALANDO PACOTES RPM", cmdsDownloadAndInstallRpm)

	cmdsDownloadFlatpaks := []string{}
	for _, packageId := range PackagesFlatpak {
		cmd := fmt.Sprintf("flatpak install -y flathub %s", packageId)
		cmdsDownloadFlatpaks = append(cmdsDownloadFlatpaks, cmd)
	}
	execCommandsGroup("INSTALANDO PACOTES FLATPAK", cmdsDownloadFlatpaks)

	toolboxVersion := "jetbrains-toolbox-2.6.2.41321"
	execCommandsGroup("INSTALANDO JETBRAINS TOOLBOX", []string{
		fmt.Sprintf("wget https://download-cdn.jetbrains.com/toolbox/%s.tar.gz -O %s/toolbox.tar.gz", toolboxVersion, tmpDir),
		fmt.Sprintf("tar -xf %s/toolbox.tar.gz -C %s", tmpDir, tmpDir),
		fmt.Sprintf("mkdir -p %s/%s", homeDir, AppImageDir),
		fmt.Sprintf("mv %s/%s/jetbrains-toolbox %s/%s", tmpDir, toolboxVersion, homeDir, AppImageDir),
		fmt.Sprintf("rm -rf %s/toolbox.tar.gz %s/%s", tmpDir, tmpDir, toolboxVersion),
	})

	execCommandsGroup("Instalando ZSH", []string{
		fmt.Sprintf("wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh -O %s/install.sh", tmpDir),
		fmt.Sprintf("chmod +x %s/install.sh", tmpDir),
		fmt.Sprintf("%s/install.sh --unattended", tmpDir),
		fmt.Sprintf("rm %s/install.sh", tmpDir),
	})

	printGroupName("CONFIGURANDO GIT")
	gitConfigPath := fmt.Sprintf("%s/.gitconfig", homeDir)
	err := os.RemoveAll(gitConfigPath)
	cobra.CheckErr(err)
	err = os.WriteFile(gitConfigPath, gitConfigData, 0744)
	cobra.CheckErr(err)

	printGroupName("CONFIGURANDO SSH")
	sshConfigPath := fmt.Sprintf("%s/.ssh/config", homeDir)
	err = os.RemoveAll(sshConfigPath)
	cobra.CheckErr(err)
	err = os.MkdirAll(filepath.Dir(sshConfigPath), 0744)
	cobra.CheckErr(err)
	err = os.WriteFile(sshConfigPath, sshConfigData, 0744)
	cobra.CheckErr(err)

	printGroupName("CONFIGURANDO 1Password")
	onePasswordConfigPath := fmt.Sprintf("%s/.config/1Password/ssh/agent.toml", homeDir)
	err = os.RemoveAll(sshConfigPath)
	cobra.CheckErr(err)
	err = os.MkdirAll(filepath.Dir(onePasswordConfigPath), 0744)
	cobra.CheckErr(err)
	err = os.WriteFile(onePasswordConfigPath, onePasswordConfigData, 0744)
	cobra.CheckErr(err)

	printGroupName("CONFIGURANDO GNOME")
	printGroupCommand("Configurando Tema Escuro")
	err = terminal.RunCommandRealtime("dconf write /org/gnome/desktop/interface/color-scheme \"'prefer-dark'\"", terminalOptions)
	cobra.CheckErr(err)

	printGroupCommand("Habilitando workspaces em todos os monitores")
	err = terminal.RunCommandRealtime("dconf write /org/gnome/mutter/workspaces-only-on-primary \"false\"", terminalOptions)
	cobra.CheckErr(err)

	printGroupCommand("Habilitando timezone automático")
	err = terminal.RunCommandRealtime("dconf write /org/gnome/desktop/datetime/automatic-timezone \"true\"", terminalOptions)
	cobra.CheckErr(err)

	printGroupCommand("Configurando formato de horas para 24h")
	err = terminal.RunCommandRealtime("dconf write /org/gnome/desktop/interface/clock-format \"'24h'\"", terminalOptions)
	cobra.CheckErr(err)
	err = terminal.RunCommandRealtime("dconf write /org/gtk/settings/file-chooser/clock-format \"'24h'\"", terminalOptions)
	cobra.CheckErr(err)

	printGroupCommand("Habilitando updates automáticos da gnome-software")
	err = terminal.RunCommandRealtime("dconf write /org/gnome/software/download-updates \"true\"", terminalOptions)
	cobra.CheckErr(err)

	printGroupCommand("Desabilitando restaurar sessão no terminal")
	err = terminal.RunCommandRealtime("dconf write /org/gnome/Ptyxis/restore-session \"false\"", terminalOptions)
	cobra.CheckErr(err)

	printGroupCommand("Configurando layouts de teclado")
	err = terminal.RunCommandRealtime("dconf write /org/gnome/desktop/input-sources/sources \"[('xkb', 'br'), ('xkb', 'us+intl')]\"", terminalOptions)
	cobra.CheckErr(err)

	printGroupCommand("Configurando velocidade do touchpad")
	err = terminal.RunCommandRealtime("dconf write /org/gnome/desktop/peripherals/touchpad/speed \"0.20171673819742497\"", terminalOptions)
	cobra.CheckErr(err)

	execCommandsGroup("Instalando tema ZSH: Starship", []string{
		fmt.Sprintf("wget https://starship.rs/install.sh -O %s/install.sh", tmpDir),
		fmt.Sprintf("chmod +x %s/install.sh", tmpDir),
		fmt.Sprintf("%s/install.sh", tmpDir),
		fmt.Sprintf("rm %s/install.sh", tmpDir),
		fmt.Sprintf("starship preset jetpack -o ~/.config/starship.toml"),
	})

	for _, font := range NerdFonts {
		execCommandsGroup(fmt.Sprintf("Instalando NerdFont : %s", font), []string{
			fmt.Sprintf("wget https://github.com/ryanoasis/nerd-fonts/releases/download/v3.4.0/%s.zip", font),
			fmt.Sprintf("unzip \"%s/%s.zip\" -d %s/%s", tmpDir, font, tmpDir, font),
			fmt.Sprintf("mkdir -p %s/.local/share/fonts", homeDir),
			fmt.Sprintf("mv %s/%s/*.ttf %s/.local/share/fonts", tmpDir, font, homeDir),
			fmt.Sprintf("rm -rf %s/%s*", tmpDir, font),
		})
	}

	printGroupCommand("Atualizando cache das fontes")
	err = terminal.RunCommandRealtime("fc-cache", terminalOptions)
	cobra.CheckErr(err)

	execCommandsGroup("Setar ZSH como shell padrão", []string{
		fmt.Sprintf("chsh -s /bin/zsh"),
	})
}
