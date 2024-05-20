#! /bin/bash
function create_config(){
local conf_dirs=("nvim" "dbus" "alacritty" "smth_random")
local config="/root/.config"
if [ ! -f $config ] 
then
    mkdir $config
    else
        echo "file doesn't work"
fi

for i in "${conf_dirs[@]}"
do
    local joined="$config/$i"
    mkdir  "$joined"
    echo "$joined was created"
done
   
}



create_config
