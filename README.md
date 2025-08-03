# Wonkyzonker

Wonkyzonker wonkily zonks your files from one directory to another, sorting them by file type and creation(modification) date.

Particularly useful for exporting your photos from your camera storage to your computer.

It tries its best not to duplicate files by checking if a file with the same name and SHA256 already exists. 
If there's a conflict, it appends a suffix `-{indx}` (where `indx` is an Nbit int [depending on your OS]) to the file name and tries again; this step is repeated until it succeeds (overflowing `indx` isn't tested)
