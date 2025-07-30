# namefix
Change names of files to match my preference - this is integrated into aviary now

Removes things according to the config

Always replaces space and . with - 
might add feature to config to have "word separator" field

Regex to match repeat characters after running the replacements is hardcoded in the fixname func for now

TODO:

flags

    config flag to open and edit config in a text editor (endgame would be built in config add/remover but the more complex the config the more complex the editor)
    specify config if needed
        possibly have update config feature to copy a given config into the config dir
    specify directory rename (currently dir is skipped)

features

    only show changed files
    only actually change files with a change to perform
    move to checking filename for change -> create a slice of structs with name and changed name to operate on

config

    specify word separator (, . - _ etc)
    contain regex??
    stand alone conf package in project? half the code is dealing with config file
    
move to more detailed config

    thinking sections in deterministic order:
    remove - anything to be replaced with empty string, can be a slice of strings in conf
    replace - anything replaced with not empty string - must be the struct we have now
    clean - clean up repeat patterns etc - ?use regex to match the characters specified in the conf?
    special - replace special characters with word separator
    finalise - re-apply file extension and change -. to . (- would be the word separator specified)
    
    

    
