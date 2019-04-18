#!/bin/bash
set -euo pipefail

WORKSPACE_PATH=~/Workspace

pp() {
  echo "====> EXECUTING STEP: $1"
}

wait_for_confirmation() {
  read -p "Have all applications been manually configured? [Yy]" -n 1 -r
  echo ""
  if [[ ! $REPLY =~ ^[Yy]$ ]]
  then
      exit 1
  fi
}
create_workspace_folder() {
  mkdir -p $WORKSPACE_PATH
}

clone_app() {
  APP=$1

  if [ -d "$WORKSPACE_PATH/$APP" ]
  then
    cd "$WORKSPACE_PATH/$APP"
    git pull --rebase
    return 0
  fi

  cd $WORKSPACE_PATH
  git clone "https://github.com/jmpak/$APP.git"
}

configure_mvim() {
  ln -sfn $WORKSPACE_PATH/dotbash/inputrc ~/.inputrc
}

configure_dotbash() {
  ln -sfn $WORKSPACE_PATH/dotbash/inputrc ~/.inputrc
}

install_brew() {
  if hash brew 2>/dev/null
  then
    return 0
  fi

  /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
  brew tap Homebrew/bundle
}

install_yadr() {
  sh -c "`curl -fsSL https://raw.githubusercontent.com/skwp/dotfiles/master/install.sh `"
}

install_brew_apps() {
  cd $WORKSPACE_PATH/dotbash
  brew bundle
}

install_pip_apps() {
  cd $WORKSPACE_PATH/dotbash
  pip3 install --user -r requirements.txt
}

manually_configure_apps() {
  echo "The below applications will require manual login:"
  echo "- Dropbox (login)"
  echo "- Google Drive (login, setup shortcuts)"
  echo "- Google Chrome (login)"
  echo "- Evernote (login)"
  echo "- Spotify (login)"
  echo "- Skype (login)"
  echo "- NordVPN (login)"
  echo "- Tunnelblick (login)"
  echo "- 1Password 6"
  echo "Docker"
  echo "VirtualBox"

  echo "The below applications will require manual setup:"
  echo "- MindNode (license)"
  echo "- iTerm (config in dotbash)"
  echo "- Adobe Creative Cloud (login, Photoshop, Lightroom)"
  echo "- Amphetamine (app store, config)"
}

configure_dotvim() {
  ln -sfn $WORKSPACE_PATH/dotvim/gvimrc ~/.gvimrc
  ln -sfn $WORKSPACE_PATH/dotvim/vimrc ~/.vimrc
  ln -sfn $WORKSPACE_PATH/dotvim/vim ~/.vim
  ln -sfn $WORKSPACE_PATH/dotvim/ctags ~/.ctags

  if [ -d "$WORKSPACE_PATH/powerline-fonts" ]
  then
    cd $WORKSPACE_PATH/powerline-fonts
    git pull --rebase
  else
    git clone https://github.com/powerline/fonts.git "$WORKSPACE_PATH/powerline-fonts"
  fi
  cd $WORKSPACE_PATH/powerline-fonts && ./install.sh
}

configure_dotgit() {
  ln -sfn $WORKSPACE_PATH/dotgit/gitconfig ~/.gitconfig
  ln -sfn $WORKSPACE_PATH/dotgit/gitignore ~/.gitignore
}

update_all_apps() {
  $WORKSPACE_PATH/dotbash/bash/bin/updateall
}

refresh_all_apps() {
  ./refresh
}

setup_ruby() {
  eval "$(rbenv init -)"
  if [ ! -d ~/.rbenv/plugins/rbenv-bundle-exec ]
  then
    git clone https://github.com/maljub01/rbenv-bundle-exec.git ~/.rbenv/plugins/rbenv-bundle-exec
  fi
  rbenv install 2.3.1 --skip-existing
  rbenv global 2.3.1
}

prereq_apps() {
  echo "The below applications will require manual install:"
  echo "- Xcode"
  wait_for_confirmation
}

pp "Creating workspace folder"      && create_workspace_folder
pp "Prerequisites"                  && prereq_apps
pp "Installing brew"                && install_brew
pp "Downloading dotbash"            && clone_app dotbash
pp "Installing brew applications"   && install_brew_apps
pp "Installing pip applications"    && install_pip_apps
pp "Configuring dotbash"            && configure_dotbash
pp "Configuring mvim"                && configure_mvim
pp "Setup Ruby"                     && setup_ruby
pp "Installing yadr"                && install_yadr
pp "Update all applications"        && update_all_apps
pp "Refresh packages"               && refresh_all_apps
pp "Manually configure apps"        && manually_configure_apps
