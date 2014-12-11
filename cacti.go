package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type cactiDS struct {
	Data_source_path string
}

func getCactiRRDs(dbCacti *sqlx.DB, graphTreeItemsIDs []string) (dataSources []cactiDS, err error) {
	start := time.Now()
	query := fmt.Sprintf("SELECT order_key FROM graph_tree_items WHERE id IN ( %s )", strings.Join(graphTreeItemsIDs, ","))

	log.Println(query)

	var OrderKeys []string
	err = dbCacti.Select(&OrderKeys, query)
	if err != nil {
		return
	}

	// I don't actually understand what order_key is, I assume
	// it's a mask of some kind. So I'll attempt to treat it
	// as such by remove trailing zeros.
	OrderKeyMasks := make([]string, len(OrderKeys))
	for k, OrderKey := range OrderKeys {
		OrderKeyMasks[k] = strings.TrimRight(OrderKey, "0")
	}

	log.Printf("Masks: %#v", OrderKeyMasks)

	/* Get non host graph entries where the graph entry is not a gprint */
	graphQuery := fmt.Sprintf(`SELECT data_source_path
  FROM graph_tree_items gtis
  JOIN graph_templates_item gti ON ( gtis.local_graph_id = gti.local_graph_id )
  JOIN data_template_rrd dtr ON ( gti.task_item_id = dtr.id )
  JOIN data_local dl ON ( dtr.local_data_id = dl.id )
  JOIN data_template_data dtd ON ( dl.id = dtd.local_data_id )
 WHERE gtis.order_key REGEXP '^(%s)'
       AND graph_type_id != 9`, strings.Join(OrderKeyMasks, "|"))

	log.Println(graphQuery)

	/* Get host graphs */
	hostQuery := fmt.Sprintf(`SELECT data_source_path
  FROM graph_tree_items gtis
  JOIN data_local dl ON ( gtis.host_id = dl.host_id )
  JOIN data_template_data dtd ON ( dl.id = dtd.local_data_id )
 WHERE gtis.order_key REGEXP '^(%s)'
   AND gtis.host_id != 0`, strings.Join(OrderKeyMasks, "|"))

	log.Println(hostQuery)

	err = dbCacti.Select(&dataSources, fmt.Sprintf("%s UNION %s", graphQuery, hostQuery))
	if err != nil {
		return
	}

	log.Println("Fetched RRD information from DB in:", time.Since(start))

	return

}
