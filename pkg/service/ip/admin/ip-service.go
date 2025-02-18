package admin

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/metal-stack/api-server/pkg/db/repository"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	"github.com/metal-stack/api/go/metalstack/admin/v2/adminv2connect"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/tag"
)

type Config struct {
	Log  *slog.Logger
	Repo *repository.Repostore
}

type ipServiceServer struct {
	log  *slog.Logger
	repo *repository.Repostore
}

func New(c Config) adminv2connect.IPServiceHandler {
	return &ipServiceServer{
		log:  c.Log.WithGroup("adminIpService"),
		repo: c.Repo,
	}
}

func (i *ipServiceServer) List(ctx context.Context, rq *connect.Request[adminv2.IPServiceListRequest]) (*connect.Response[adminv2.IPServiceListResponse], error) {
	i.log.Debug("list", "ip", rq)
	req := rq.Msg

	adminReq := &apiv2.IPServiceListRequest{
		Ip:               req.Ip,
		Network:          req.Network,
		Project:          req.Project,
		Name:             req.Name,
		Uuid:             req.Uuid,
		MachineId:        req.MachineId,
		ParentPrefixCidr: req.ParentPrefixCidr,
		Tags:             req.Tags,
		Type:             req.Type,
		AddressFamily:    req.AddressFamily,
	}

	resp, err := i.repo.IP(nil).List(ctx, adminReq)
	if err != nil {
		return nil, err
	}

	var res []*apiv2.IP
	for _, ip := range resp {

		m := tag.NewTagMap(ip.Tags)
		if _, ok := m.Value(tag.MachineID); ok {
			// we do not want to show machine ips (e.g. firewall public ips)
			continue
		}

		converted, err := i.repo.IP(nil).ConvertToProto(ip)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		res = append(res, converted)
	}

	return connect.NewResponse(&adminv2.IPServiceListResponse{
		Ips: res,
	}), nil
}

func (i *ipServiceServer) Issues(ctx context.Context, rq *connect.Request[adminv2.IPServiceIssuesRequest]) (*connect.Response[adminv2.IPServiceIssuesResponse], error) {
	panic("unimplemented")
}
