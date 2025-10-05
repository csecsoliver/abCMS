pip install pex
pex . --entry-point main -o app.pex --scie lazy \
  --include-file *.html \
  --include-file ../frontend:frontend