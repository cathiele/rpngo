package tinyfs

// This doesn't handle more complex cases like /foo/bar/../baz
// Maybe it can be reimplemented later.
//
// TODO: tralingSlash is always false so far...
func absPath(pwd string, path string, leadingSlash bool, trailingSlash bool) string {
	if path == "." {
		path = ""
	} else if path == ".." {
		path = parentOfPwd(pwd)
	}
	if (len(path) == 0) || (path[0] != '/') {
		path = pwd + "/" + path
	}
	if path == "/" {
		if leadingSlash || trailingSlash {
			return "/"
		}
		return ""
	}
	if leadingSlash {
		path = addLeadingSlash(path)
	} else {
		path = removeLeadingSlash(path)
	}
	if trailingSlash {
		path = addTrailingSlash(path)
	} else {
		path = removeTrailingSlash(path)
	}
	if path == "//" {
		path = "/"
	}
	return path
}

func basePath(p string) string {
	slashIdx := -1
	for i, c := range p {
		if c == '/' {
			slashIdx = i
		}
	}
	if slashIdx < 0 {
		return p
	}
	return p[slashIdx+1:]
}

func addLeadingSlash(p string) string {
	if p[0] == '/' {
		return p
	}
	return "/" + p
}

func removeLeadingSlash(p string) string {
	if p[0] == '/' {
		return p[1:]
	}
	return p
}

func addTrailingSlash(p string) string {
	if p[len(p)-1] == '/' {
		return p
	}
	return p + "/"
}

func removeTrailingSlash(p string) string {
	if p[len(p)-1] == '/' {
		return p[:len(p)-1]
	}
	return p
}

func parentOfPwd(pwd string) string {
	lastSlashIdx := 0
	for i, c := range pwd {
		if c == '/' {
			lastSlashIdx = i
		}
	}
	if lastSlashIdx == 0 {
		return "/"
	}
	return pwd[:lastSlashIdx]
}
