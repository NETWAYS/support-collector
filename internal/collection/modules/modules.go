package modules

import (
	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/collection/modules/ansible"
	"github.com/NETWAYS/support-collector/internal/collection/modules/base"
	"github.com/NETWAYS/support-collector/internal/collection/modules/corosync"
	"github.com/NETWAYS/support-collector/internal/collection/modules/elastic"
	"github.com/NETWAYS/support-collector/internal/collection/modules/foreman"
	"github.com/NETWAYS/support-collector/internal/collection/modules/grafana"
	"github.com/NETWAYS/support-collector/internal/collection/modules/graphite"
	"github.com/NETWAYS/support-collector/internal/collection/modules/graylog"
	"github.com/NETWAYS/support-collector/internal/collection/modules/icinga2"
	"github.com/NETWAYS/support-collector/internal/collection/modules/icingadb"
	"github.com/NETWAYS/support-collector/internal/collection/modules/icingadirector"
	"github.com/NETWAYS/support-collector/internal/collection/modules/icingaweb2"
	"github.com/NETWAYS/support-collector/internal/collection/modules/influxdb"
	"github.com/NETWAYS/support-collector/internal/collection/modules/keepalived"
	"github.com/NETWAYS/support-collector/internal/collection/modules/mongodb"
	"github.com/NETWAYS/support-collector/internal/collection/modules/mysql"
	"github.com/NETWAYS/support-collector/internal/collection/modules/postgresql"
	"github.com/NETWAYS/support-collector/internal/collection/modules/prometheus"
	"github.com/NETWAYS/support-collector/internal/collection/modules/puppet"
	"github.com/NETWAYS/support-collector/internal/collection/modules/webservers"
)

var (
	List = map[string]func(*collection.Collection){
		"base":            base.CollectLocal,
		"webservers":      webservers.Collect,
		"icinga2":         icinga2.Collect,
		"icingaweb2":      icingaweb2.Collect,
		"icinga-director": icingadirector.Collect,
		"elastic":         elastic.Collect,
		"corosync":        corosync.Collect,
		"keepalived":      keepalived.Collect,
		"mongodb":         mongodb.Collect,
		"mysql":           mysql.Collect,
		"influxdb":        influxdb.Collect,
		"postgresql":      postgresql.Collect,
		"prometheus":      prometheus.Collect,
		"ansible":         ansible.Collect,
		"puppet":          puppet.Collect,
		"grafana":         grafana.Collect,
		"graphite":        graphite.Collect,
		"graylog":         graylog.Collect,
		"icingadb":        icingadb.Collect,
		"foreman":         foreman.Collect,
	}

	Order = []string{
		"base",
		//"webservers",
		//"icinga2",
		//"icingaweb2",
		//"icinga-director",
		//"icingadb",
		//"elastic",
		//"corosync",
		//"keepalived",
		//"mongodb",
		//"mysql",
		//"influxdb",
		//"postgresql",
		//"prometheus",
		//"ansible",
		//"puppet",
		//"grafana",
		//"graphite",
		//"graylog",
		//"foreman",
	} // TODO enable modules
)
