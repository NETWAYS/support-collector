package modules

import (
	"github.com/NETWAYS/support-collector/modules/ansible"
	"github.com/NETWAYS/support-collector/modules/base"
	"github.com/NETWAYS/support-collector/modules/corosync"
	"github.com/NETWAYS/support-collector/modules/elastic"
	"github.com/NETWAYS/support-collector/modules/foreman"
	"github.com/NETWAYS/support-collector/modules/grafana"
	"github.com/NETWAYS/support-collector/modules/graphite"
	"github.com/NETWAYS/support-collector/modules/graylog"
	"github.com/NETWAYS/support-collector/modules/icinga2"
	"github.com/NETWAYS/support-collector/modules/icingadb"
	"github.com/NETWAYS/support-collector/modules/icingadirector"
	"github.com/NETWAYS/support-collector/modules/icingaweb2"
	"github.com/NETWAYS/support-collector/modules/influxdb"
	"github.com/NETWAYS/support-collector/modules/keepalived"
	"github.com/NETWAYS/support-collector/modules/mongodb"
	"github.com/NETWAYS/support-collector/modules/mysql"
	"github.com/NETWAYS/support-collector/modules/postgresql"
	"github.com/NETWAYS/support-collector/modules/prometheus"
	"github.com/NETWAYS/support-collector/modules/puppet"
	"github.com/NETWAYS/support-collector/modules/webservers"
	"github.com/NETWAYS/support-collector/pkg/collection"
)

var (
	List = map[string]func(*collection.Collection){
		"base":            base.Collect,
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
