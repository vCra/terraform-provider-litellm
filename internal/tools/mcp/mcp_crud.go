package mcp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/scalepad/terraform-provider-litellm/internal/litellm"
)

func createMCPServer(ctx context.Context, c *litellm.Client, server *MCPServer) (*MCPServerResponse, error) {
	resp, err := c.SendRequest(ctx, http.MethodPost, "/v1/mcp/server", server)
	if err != nil {
		return nil, err
	}

	return parseMCPServerAPIResponse(resp)
}

func getMCPServer(ctx context.Context, c *litellm.Client, serverID string) (*MCPServerResponse, error) {
	resp, err := c.SendRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/mcp/server/%s", serverID), nil)
	if err != nil {
		// Check if it's a not found error
		if litellm.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return parseMCPServerAPIResponse(resp)
}

func updateMCPServer(ctx context.Context, c *litellm.Client, server *MCPServer) (*MCPServerResponse, error) {
	// Create a new map with only the fields that can be updated
	updateData := map[string]interface{}{
		"server_id":         server.ServerID,
		"server_name":       server.ServerName,
		"alias":             server.Alias,
		"description":       server.Description,
		"url":               server.URL,
		"transport":         server.Transport,
		"spec_version":      server.SpecVersion,
		"auth_type":         server.AuthType,
		"mcp_info":          server.MCPInfo,
		"mcp_access_groups": server.MCPAccessGroups,
		"command":           server.Command,
		"args":              server.Args,
		"env":               server.Env,
	}

	resp, err := c.SendRequest(ctx, http.MethodPut, "/v1/mcp/server", updateData)
	if err != nil {
		// If server not found during update, return error
		if litellm.IsNotFound(err) {
			return nil, err
		}
		return nil, err
	}

	return parseMCPServerAPIResponse(resp)
}

func deleteMCPServer(ctx context.Context, c *litellm.Client, serverID string) error {
	_, err := c.SendRequest(ctx, http.MethodDelete, fmt.Sprintf("/v1/mcp/server/%s", serverID), nil)
	// If it's a not found error, consider it successful (already deleted)
	if err != nil && litellm.IsNotFound(err) {
		return nil
	}
	return err
}
