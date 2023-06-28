/**
 * Copyright 2023 Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package pool

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"

	v1pools "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/pools/v1"
)

func ListPools(w http.ResponseWriter, r *http.Request) {

	// TODO: This needs to page for the end client - iterator blasts through everything

	req := &v1pools.ListPoolsRequest{}

	log.Debugf("ListPools request: %v", req)
	iter := waas.GetClients().PoolService.ListPools(r.Context(), req)

	var pools []*v1pools.Pool
	for {
		pool, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Errorf("cannot iterate pools: %v", err)
			utils.HttpBadGateway(w)
			return
		}

		pools = append(pools, pool)
	}

	response := &v1pools.ListPoolsResponse{Pools: pools}

	log.Debugf("ListPools response: %v", response)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("cannot marshal and write list pools response: %v", err)
		utils.HttpBadGateway(w)
	}
}
