package service

// type transferProxy struct {
// 	get  kitendpoint.Endpoint
// 	list kitendpoint.Endpoint
// }

// func (proxy transferProxy) Get(ip, user, passwd, remoteFilePath, localPath string) (string, error) {
// 	resp, err := proxy.get(nil, endpoint.GetRequest{
// 		Ip:             ip,
// 		User:           user,
// 		Passwd:         passwd,
// 		RemoteFilePath: remoteFilePath,
// 		LocalPath:      localPath,
// 	})
// 	if err != nil {
// 		return "", err
// 	}
// 	return resp.(endpoint.GetResponse).V, nil
// }

// func (proxy transferProxy) List(ip, user, passwd, remoteFilePath string) ([]string, error) {
// 	resp, err := proxy.list(nil, endpoint.ListRequest{
// 		Ip:             ip,
// 		User:           user,
// 		Passwd:         passwd,
// 		RemoteFilePath: remoteFilePath,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp.(endpoint.ListResponse).V, nil
// }

// func MakeGetProxy(_url string) (kitendpoint.Endpoint, error) {
// 	u, err := url.Parse(_url)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return httptransport.NewClient(
// 		"POST",
// 		u,
// 		encodeRequest,
// 		decodeGetResponse,
// 	).Endpoint(), nil
// }

// func MakeListProxy(_url string) (kitendpoint.Endpoint, error) {
// 	u, err := url.Parse(_url)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return httptransport.NewClient(
// 		"POST",
// 		u,
// 		encodeRequest,
// 		decodeListResponse,
// 	).Endpoint(), nil
// }
