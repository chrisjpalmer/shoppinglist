live:
	ENTR_INOTIFY_WORKAROUND=y find ./ ! -name '*templ.go' | entr -r make serve

serve:
	go tool templ generate
	go run .
