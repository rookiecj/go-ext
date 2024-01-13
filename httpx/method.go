package httpx

//
//func fetch[Body any](client *Client, method string, url *url.URL, options ...Option) (body Body, err error) {
//	if client == nil {
//		client = DefaultClient
//	}
//
//	resp, err := client.Do(method, url, options...)
//	if err != nil {
//		err = fmt.Errorf("failed to request: %s", err)
//		return
//	}
//
//	// parse body
//	if contentTypes, ok := resp.Header["Content-Type"]; ok {
//		if bodyParser := getBodyParser(contentTypes[0]); bodyParser != nil {
//			err = bodyParser(resp.Data, &body)
//			if err == nil {
//				return
//			}
//		}
//	}
//
//	// string
//	bodyType := reflect.TypeOf(body)
//	if bodyType.Kind() == reflect.String {
//		valType := reflect.ValueOf(body)
//		valRef := reflect.Indirect(valType)
//		valRef.SetString(string(resp.Data))
//		err = nil
//		return
//	}
//
//	err = fmt.Errorf("parse error: %s", err)
//	return
//}
