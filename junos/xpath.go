package junos

//const (
//	xPathSep  = "/"
//	xpathFile = "xpath_inputs.xml"
//)

//type Xpath struct {
//	Name string `xml:"name,attr"`
//}

//type Xpaths struct {
//	XMLName struct{} `xml:"file-list"`
//	Xpaths  []Xpath  `xml:"xpath"`
//}

//// deviceConfigXpaths reads the specified XML file and returns []string
//// representing the deepest paths encountered while parsing the file.
//func deviceConfigXpaths(filename string) (Xpaths, error) {
//	var result Xpaths
//
//	f, err := os.Open(filename)
//	if err != nil {
//		return result, fmt.Errorf("while opening device config file %q - %w", filename, err)
//	}
//	defer func(closer io.Closer) { _ = closer.Close() }(f) // ignoring the error on read seems reasonable
//
//	xmlDec := xml.NewDecoder(f)
//
//	var path []string
//	xPathMap := make(map[string]struct{})
//
//	for {
//		tok, err := xmlDec.Token()
//		if err != nil {
//			if errors.Is(err, io.EOF) {
//				break
//			}
//
//			return result, fmt.Errorf("while getting xml token - %w", err)
//		}
//
//		switch tok := tok.(type) {
//		case xml.StartElement:
//			// keep track of where we are by appending to the current path
//			path = append(path, tok.Name.Local)
//			// record this position in the tree
//			xPathMap[strings.Join(path, xPathSep)] = struct{}{}
//		case xml.EndElement:
//			// keep track of where we are by trimming the current path
//			path = path[:len(path)-1]
//		}
//	}
//
//	// trim the map so that only leaf entries remain
//	for k := range xPathMap {
//		pathElems := strings.Split(k, xPathSep)
//		for len(pathElems) > 0 {
//			pathElems = pathElems[:len(pathElems)-1]
//			delete(xPathMap, strings.Join(pathElems, xPathSep))
//		}
//	}
//
//	result.Xpaths = make([]Xpath, len(xPathMap))
//	for i, xpath := range helpers.OrderedKeys(xPathMap) {
//		result.Xpaths[i] = Xpath{Name: xpath}
//	}
//
//	return result, nil
//}

//func deviceConfigXpathsToFile(in, out string) error {
//	xPaths, err := deviceConfigXpaths(in)
//	if err != nil {
//		return fmt.Errorf("while parsing device config %q - %w", in, err)
//	}
//
//	b, err := xml.MarshalIndent(xPaths, "", "  ")
//	if err != nil {
//		return fmt.Errorf("while marshaling xpaths from device configuration - %w", err)
//	}
//
//	err = os.WriteFile(xpathFile, b, 0o644)
//	if err != nil {
//		return fmt.Errorf("while writing xpath data to %q - %w", xpathFile, err)
//	}
//
//	return nil
//}
