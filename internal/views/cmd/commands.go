package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/johannessarpola/lutakkols/internal/views/messages"
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"github.com/johannessarpola/lutakkols/pkg/api/options"
	"github.com/johannessarpola/lutakkols/pkg/api/provider"
	"github.com/johannessarpola/lutakkols/pkg/logger"
)

func GetDetails(eventID string, eventURL string, provider provider.Provider, opts ...options.ProviderOption) tea.Cmd {
	return func() tea.Msg {
		logger.Log.Debugf("getting description for %s from provider", eventURL)
		eventDetails, err := provider.GetDetails(eventID, eventURL, opts...)
		if err != nil {
			return err
		}
		return messages.EventDescriptionFetched{Details: eventDetails, ProviderOptions: opts}
	}
}

func GetAscii(eventID string, imageURL string, provider provider.Provider, opts ...options.ProviderOption) tea.Cmd {
	return func() tea.Msg {
		logger.Log.Debugf("getting ascii with url %s from provider", imageURL)
		eventAscii, err := provider.GetAscii(eventID, imageURL, opts...)
		if err != nil {
			return err
		}
		return messages.EventAsciiFetched{Ascii: eventAscii.Ascii}
	}
}

func GetEvents(provider provider.Provider, opts ...options.ProviderOption) tea.Cmd {

	return func() tea.Msg {
		var (
			events *models.Events
			err    error
		)
		logger.Log.Debugf("getting events from provider")
		events, err = provider.GetEvents(opts...)

		if err != nil {
			return err
		}
		return messages.EventsFetched{Events: events.Events, Time: events.UpdatedAt}
	}
}
