package channelz

import (
	"context"
	"io"

	channelzgrpc "google.golang.org/grpc/channelz/grpc_channelz_v1"
	log "google.golang.org/grpc/grpclog"
)

// writeTopChannelsPage writes an HTML document to w containing per-channel RPC stats, including a header and a footer.
func (h *channelzHandler) writeTopChannelsPage(w io.Writer) {
	writeHeader(w, "ChannelZ Stats")
	h.writeTopChannels(w)
	writeFooter(w)
}

func writeHeader(w io.Writer, title string) {
	if err := headerTemplate.Execute(w, headerData{Title: title}); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

func writeFooter(w io.Writer) {
	if err := footerTemplate.Execute(w, nil); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}

}

// writeTopChannels writes HTML to w containing per-channel RPC stats.
//
// It includes neither a header nor footer, so you can embed this data in other pages.
func (h *channelzHandler) writeTopChannels(w io.Writer) {
	if err := topChannelsTemplate.Execute(w, h.getTopChannels()); err != nil {
		log.Errorf("channelz: executing template: %v", err)
	}
}

func (h *channelzHandler) getTopChannels() *channelzgrpc.GetTopChannelsResponse {
	client, err := h.connect()
	if err != nil {
		log.Errorf("Error creating channelz client %+v", err)
		return nil
	}
	ctx := context.Background()
	channels, err := client.GetTopChannels(ctx, &channelzgrpc.GetTopChannelsRequest{})
	if err != nil {
		log.Errorf("Error querying GetTopChannels %+v", err)
		return nil
	}
	return channels
}

const topChannelsTemplateHTML = `
<table frame=box cellspacing=0 cellpadding=2>
    <tr class="header">
		<th colspan=100 style="text-align:left">Top Channels: {{.Channel | len}}</th>
    </tr>

    <tr classs="header">
        <th>ID</th>
        <th>Name</th>
        <th>State</th>
        <th>Target</th>
        <th>Subchannels</th>
        <th>CreationTimestamp</th>
        <th>CallsStarted</th>
        <th>CallsSucceeded</th>
        <th>CallsFailed</th>
        <th>LastCallStartedTimestamp</th>
        <th>ChannelRef</th>
    </tr>
{{range .Channel}}
    <tr>
		{{ $channelID := .Ref.ChannelId }}
        <td>{{.Ref.ChannelId}}</td>
        <td><b>{{.Ref.Name}}</b></td>
        <td>{{.Data.State}}</td>
        <td>{{.Data.Target}}</td>
		<td>
			{{range .SubchannelRef}}
				<a href="{{$channelID}}/{{.SubchannelId}}">[{{.SubchannelId}}]{{.Name}}</a>
			{{end}}
		</td>
        <td>{{.Data.Trace.CreationTimestamp | timestamp}}</td>
        <td>{{.Data.CallsStarted}}</td>
        <td>{{.Data.CallsSucceeded}}</td>
        <td>{{.Data.CallsFailed}}</td>
        <td>{{.Data.LastCallStartedTimestamp | timestamp}}</td>
		<td>{{.ChannelRef}}</td>
	</tr>
    <tr classs="header">
        <th colspan=100>Events</th>
    </tr>
	<tr>
		<td>&nbsp;</td>
        <td colspan=100>
			<pre>
			{{- range .Data.Trace.Events}}
{{.Severity}} [{{.Timestamp}}]: {{.Description}}
			{{- end -}}
			</pre>
		</td>
    </tr>
{{end}}
</table>
`
