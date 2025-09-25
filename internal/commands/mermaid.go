// SPDX-License-Identifier: Apache-2.0
// internal/commands/mermaid.go
package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"maestro/internal/common"
)

// MermaidCommand implements the mermaid command
type MermaidCommand struct {
	*BaseCommand
	workflowFile    string
	sequenceDiagram bool
	flowchartTD     bool
	flowchartLR     bool
}

// NewMermaidCommand creates a new mermaid command
func NewMermaidCommand() *cobra.Command {
	mermaidCmd := &MermaidCommand{}

	cmd := &cobra.Command{
		Use:   "mermaid WORKFLOW_FILE",
		Short: "Generate mermaid diagrams from a workflow file",
		Long:  `Generate mermaid diagrams from a workflow file.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options := NewCommandOptions(cmd)

			mermaidCmd.BaseCommand = NewBaseCommand(options)
			mermaidCmd.workflowFile = args[0]

			// Get flag values
			var err error
			mermaidCmd.sequenceDiagram, err = cmd.Flags().GetBool("sequenceDiagram")
			if err != nil {
				return err
			}

			mermaidCmd.flowchartTD, err = cmd.Flags().GetBool("flowchart-td")
			if err != nil {
				return err
			}

			mermaidCmd.flowchartLR, err = cmd.Flags().GetBool("flowchart-lr")
			if err != nil {
				return err
			}

			return mermaidCmd.Run()
		},
	}

	// Add flags
	cmd.Flags().Bool("sequenceDiagram", false, "Sequence diagram mermaid")
	cmd.Flags().Bool("flowchart-td", false, "Flowchart TD (top down) mermaid")
	cmd.Flags().Bool("flowchart-lr", false, "Flowchart LR (left right) mermaid")

	return cmd
}

// Run executes the mermaid command
func (c *MermaidCommand) Run() error {
	// Parse the workflow YAML file
	workflowYaml, err := common.ParseYAML(c.workflowFile)
	if err != nil {
		c.Console().Error(fmt.Sprintf("Unable to parse workflow file: %s", err))
		return err
	}

	// Generate the mermaid diagram
	mermaid, err := c.generateMermaid(workflowYaml[0])
	if err != nil {
		c.Console().Error(fmt.Sprintf("Unable to generate mermaid for workflow: %s", err))
		return err
	}

	// Print the mermaid diagram
	if !c.IsSilent() {
		c.Console().Ok("Created mermaid for workflow\n")
	}
	c.Console().Print(mermaid + "\n")

	return nil
}

// generateMermaid generates a mermaid diagram from a workflow
func (c *MermaidCommand) generateMermaid(workflow common.YAMLDocument) (string, error) {
	// In the Python implementation, this calls workflow.to_mermaid()
	// We'll need to implement the equivalent functionality in Go

	// For now, we'll just return a dummy mermaid diagram
	var diagramType, direction string

	if c.sequenceDiagram {
		diagramType = "sequenceDiagram"
	} else if c.flowchartTD {
		diagramType = "flowchart"
		direction = "TD"
	} else if c.flowchartLR {
		diagramType = "flowchart"
		direction = "LR"
	} else {
		diagramType = "sequenceDiagram"
	}
	mermaid := NewMermaid(workflow, diagramType, direction)
	result, err := mermaid.ToMarkdown()
	if err != nil {
	   return "", err
	}

	// TODO: Implement the actual mermaid generation logic
	// This would involve:
	// 1. Parsing the workflow structure
	// 2. Generating the appropriate mermaid syntax

	return result, nil
}


// Mermaid represents a mermaid diagram generator
type Mermaid struct {
//	workflow    map[string]interface{}
	workflow    common.YAMLDocument
	kind        string
	orientation string
}

// NewMermaid creates a new Mermaid instance
//func NewMermaid(workflow map[string]interface{}, kind string, orientation string) *Mermaid {
func NewMermaid(workflow common.YAMLDocument, kind string, orientation string) *Mermaid {
	if kind == "" {
		kind = "sequenceDiagram"
	}
	if orientation == "" {
		orientation = "TD"
	}
	return &Mermaid{
		workflow:    workflow,
		kind:        kind,
		orientation: orientation,
	}
}

// ToMarkdown converts the workflow to a mermaid diagram in markdown format
func (m *Mermaid) ToMarkdown() (string, error) {
	if m.kind == "sequenceDiagram" {
		return m.toSequenceDiagram(), nil
	} else if m.kind == "flowchart" {
		return m.toFlowchart(), nil
	} else {
		return "", fmt.Errorf("invalid Mermaid kind: %s", m.kind)
	}
}

// fixAgentName replaces hyphens with underscores in agent names
func (m *Mermaid) fixAgentName(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}

// agentForStep returns the agent for a given step name
func (m *Mermaid) agentForStep(stepName string) string {
	template := m.workflow["spec"].(common.YAMLDocument)["template"].(common.YAMLDocument)
	steps, ok := template["steps"].([]interface{})
	if !ok {
		return ""
	}

	for _, s := range steps {
		step := s.(common.YAMLDocument)
		if name, ok := step["name"]; ok && name == stepName {
			if agent, ok := step["agent"]; ok {
				return agent.(string)
			}
		}
	}
	return ""
}

// sequenceParticipants returns the list of participants for a sequence diagram
func (m *Mermaid) sequenceParticipants() []string {
	template := m.workflow["spec"].(common.YAMLDocument)["template"].(common.YAMLDocument)

	// Check if agents are explicitly defined
	if agents, ok := template["agents"]; ok {
		agentsList := []string{}
		for _, agent := range agents.([]interface{}) {
			agentsList = append(agentsList, agent.(string))
		}
		return agentsList
	}

	// Otherwise, collect agents from steps
	seen := []string{}
	steps, ok := template["steps"].([]interface{})
	if !ok {
		return seen
	}

	for _, s := range steps {
		step := s.(common.YAMLDocument)
		if agent, ok := step["agent"]; ok {
			agentStr := agent.(string)

			// Skip steps with context or outputs
			if _, hasContext := step["context"]; hasContext {
				continue
			}
			if _, hasOutputs := step["outputs"]; hasOutputs {
				continue
			}

			// Add agent if not already seen
			found := false
			for _, a := range seen {
				if a == agentStr {
					found = true
					break
				}
			}
			if !found {
				seen = append(seen, agentStr)
			}
		}
	}

	return seen
}

// toSequenceDiagram generates a mermaid sequence diagram
func (m *Mermaid) toSequenceDiagram() string {
	var sb strings.Builder
	sb.WriteString("sequenceDiagram\n")

	// Add participants
	for _, agent := range m.sequenceParticipants() {
		sb.WriteString(fmt.Sprintf("participant %s\n", m.fixAgentName(agent)))
	}

	template := m.workflow["spec"].(common.YAMLDocument)["template"].(common.YAMLDocument)
	steps, ok := template["steps"].([]interface{})
	if !ok {
		steps = []interface{}{}
	}

	var agentL string
	for i, s := range steps {
		step := s.(common.YAMLDocument)

		// Skip scoring/context-only steps
		if _, hasContext := step["context"]; hasContext {
			continue
		}
		if _, hasOutputs := step["outputs"]; hasOutputs {
			continue
		}

		// Update agentL only when this step names a real agent
		if agent, ok := step["agent"]; ok {
			agentL = m.fixAgentName(agent.(string))
		}

		// Find next real agent for the arrow
		var agentR string
		for j := i + 1; j < len(steps); j++ {
			nextStep := steps[j].(common.YAMLDocument)

			if _, hasContext := nextStep["context"]; hasContext {
				continue
			}
			if _, hasOutputs := nextStep["outputs"]; hasOutputs {
				continue
			}

			if agent, ok := nextStep["agent"]; ok {
				agentR = m.fixAgentName(agent.(string))
				break
			}
		}

		stepName := step["name"].(string)
		if agentR != "" {
			sb.WriteString(fmt.Sprintf("%s->>%s: %s\n", agentL, agentR, stepName))
		} else {
			sb.WriteString(fmt.Sprintf("%s->>%s: %s\n", agentL, agentL, stepName))
		}

		// Handle condition / parallel / loop
		if condition, ok := step["condition"]; ok {
			conditions := condition.([]interface{})
			for _, c := range conditions {
				sb.WriteString(m.toSequenceDiagramCondition(agentL, agentR, c.(common.YAMLDocument)))
			}
		}

		if _, ok := step["parallel"]; ok {
			sb.WriteString(m.toSequenceDiagramParallel(agentL, step))
		}

		if loop, ok := step["loop"]; ok {
			sb.WriteString(m.toSequenceDiagramLoop(agentL, loop.(common.YAMLDocument)))
		}
	}

	// Global cron-event block
	if event, ok := template["event"]; ok {
		eventMap := event.(common.YAMLDocument)
		if _, hasCron := eventMap["cron"]; hasCron {
			sb.WriteString(m.toSequenceDiagramEvent(eventMap))
		}
	}

	// Global exception block
	if exc, ok := template["exception"]; ok {
		sb.WriteString(m.toSequenceDiagramException(steps, exc.(common.YAMLDocument)))
	}

	return sb.String()
}

// toSequenceDiagramEvent generates the event part of a sequence diagram
func (m *Mermaid) toSequenceDiagramEvent(event common.YAMLDocument) string {
	var sb strings.Builder
	name, _ := event["name"].(string)
	cron, _ := event["cron"].(string)
	exit, _ := event["exit"].(string)

	sb.WriteString(fmt.Sprintf("alt cron \"%s\"\n", cron))

	if steps, ok := event["steps"]; ok {
		for _, stepName := range steps.([]interface{}) {
			agent := m.agentForStep(stepName.(string))
			sb.WriteString(fmt.Sprintf("  cron->>%s: %s\n", agent, stepName))
		}
	} else {
		agent, _ := event["agent"].(string)
		sb.WriteString(fmt.Sprintf("  cron->>%s: %s\n", agent, name))
	}

	sb.WriteString("else\n")
	sb.WriteString(fmt.Sprintf("  cron->>exit: %s\n", exit))
	sb.WriteString("end\n")

	return sb.String()
}

// toSequenceDiagramParallel generates the parallel part of a sequence diagram
func (m *Mermaid) toSequenceDiagramParallel(agentL string, parallelStep common.YAMLDocument) string {
	var sb strings.Builder
	sb.WriteString("par\n")

	parallel := parallelStep["parallel"].([]interface{})
	for i, agent := range parallel {
		agentR := m.fixAgentName(agent.(string))
		sb.WriteString(fmt.Sprintf("  %s->>%s: %s\n", agentL, agentR, parallelStep["name"]))

		if i < len(parallel)-1 {
			sb.WriteString("and\n")
		}
	}

	sb.WriteString("end\n")
	return sb.String()
}

// toSequenceDiagramLoop generates the loop part of a sequence diagram
func (m *Mermaid) toSequenceDiagramLoop(agentL string, loopDef common.YAMLDocument) string {
	var sb strings.Builder
	expr := "True"

	if until, ok := loopDef["until"]; ok {
		expr = until.(string)
	}

	sb.WriteString(fmt.Sprintf("loop %s\n", expr))

	agent, _ := loopDef["agent"].(string)
	loopType := "until"
	if _, ok := loopDef["until"]; !ok {
		loopType = "loop"
	}

	sb.WriteString(fmt.Sprintf("  %s-->%s: %s\n", agentL, m.fixAgentName(agent), loopType))
	sb.WriteString("end\n")

	return sb.String()
}

// toSequenceDiagramCondition generates the condition part of a sequence diagram
func (m *Mermaid) toSequenceDiagramCondition(agentL string, agentR string, condition common.YAMLDocument) string {
	var sb strings.Builder

	if caseVal, ok := condition["case"]; ok {
		cond := caseVal.(string)
		do := ""

		if doVal, ok := condition["do"]; ok {
			do = doVal.(string)
		}

		if _, ok := condition["default"]; ok {
			cond = "default"
			do = condition["default"].(string)
		}

		sb.WriteString(fmt.Sprintf("%s->>%s: %s %s\n", agentL, agentR, do, cond))
	} else if ifVal, ok := condition["if"]; ok {
		ifExpr := ifVal.(string)
		thenExpr := ""

		if thenVal, ok := condition["then"]; ok {
			thenExpr = thenVal.(string)
		}

		sb.WriteString(fmt.Sprintf("%s->>%s: %s\n", agentL, agentR, ifExpr))
		sb.WriteString("alt if True\n")
		sb.WriteString(fmt.Sprintf("  %s->>%s: %s\n", agentL, agentR, thenExpr))

		if elseVal, ok := condition["else"]; ok {
			elseExpr := elseVal.(string)
			sb.WriteString("else is False\n")
			sb.WriteString(fmt.Sprintf("  %s->>%s: %s\n", agentR, agentL, elseExpr))
		}

		sb.WriteString("end\n")
	}

	return sb.String()
}

// toSequenceDiagramException generates the exception part of a sequence diagram
func (m *Mermaid) toSequenceDiagramException(steps []interface{}, exception common.YAMLDocument) string {
	var sb strings.Builder
	sb.WriteString("alt exception\n")

	for _, s := range steps {
		step := s.(common.YAMLDocument)
		if agent, ok := step["agent"]; ok {
			agentL := m.fixAgentName(agent.(string))
			exceptionAgent := exception["agent"].(string)
			exceptionName := exception["name"].(string)
			sb.WriteString(fmt.Sprintf("  %s->>%s: %s\n", agentL, exceptionAgent, exceptionName))
		}
	}

	sb.WriteString("end")
	return sb.String()
}

// toFlowchart generates a mermaid flowchart
func (m *Mermaid) toFlowchart() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("flowchart %s\n", m.orientation))

	template := m.workflow["spec"].(common.YAMLDocument)["template"].(common.YAMLDocument)
	steps, ok := template["steps"].([]interface{})
	if !ok {
		steps = []interface{}{}
	}

	i := 0
	for i < len(steps) {
		step := steps[i].(common.YAMLDocument)

		// Skip scoring/context-only steps
		if _, hasContext := step["context"]; hasContext {
			i++
			continue
		}
		if _, hasOutputs := step["outputs"]; hasOutputs {
			i++
			continue
		}

		aL, _ := step["agent"].(string)

		// Find next real step
		var aR string
		for j := i + 1; j < len(steps); j++ {
			nextStep := steps[j].(common.YAMLDocument)

			if _, hasContext := nextStep["context"]; hasContext {
				continue
			}
			if _, hasOutputs := nextStep["outputs"]; hasOutputs {
				continue
			}

			if agent, ok := nextStep["agent"]; ok {
				aR = agent.(string)
				break
			}
		}

		stepName := step["name"].(string)
		if aR != "" {
			sb.WriteString(fmt.Sprintf("%s-- %s -->%s\n", aL, stepName, aR))
		} else {
			sb.WriteString(fmt.Sprintf("%s-- %s -->%s\n", aL, stepName, aL))
		}

		if condition, ok := step["condition"]; ok {
			conditions := condition.([]interface{})
			for _, c := range conditions {
				sb.WriteString(m.toFlowchartCondition(aL, aR, step, c.(common.YAMLDocument)))
			}
		}

		i++
	}

	// Global exception block
	if exc, ok := template["exception"]; ok {
		sb.WriteString(m.toFlowchartException(steps, exc.(common.YAMLDocument)))
	}

	return sb.String()
}

// toFlowchartCondition generates the condition part of a flowchart
func (m *Mermaid) toFlowchartCondition(agentL string, agentR string, step common.YAMLDocument, condition common.YAMLDocument) string {
	var sb strings.Builder

	if caseVal, ok := condition["case"]; ok {
		cond := caseVal.(string)
		do := ""

		if doVal, ok := condition["do"]; ok {
			do = doVal.(string)
		}

		if _, ok := condition["default"]; ok {
			cond = "default"
			do = condition["default"].(string)
		}

		sb.WriteString(fmt.Sprintf("%s-- %s %s -->%s\n", agentL, do, cond, agentR))
	}

	if ifVal, ok := condition["if"]; ok {
		expr := ifVal.(string)
		thenExpr := ""
		elseExpr := ""

		if thenVal, ok := condition["then"]; ok {
			thenExpr = thenVal.(string)
		}

		if elseVal, ok := condition["else"]; ok {
			elseExpr = elseVal.(string)
		}

		stepName := step["name"].(string)
		sb.WriteString(fmt.Sprintf("%s --> Condition{\"%s\"}\n", stepName, expr))
		sb.WriteString(fmt.Sprintf("  Condition -- Yes --> %s\n", thenExpr))
		sb.WriteString(fmt.Sprintf("  Condition -- No --> %s\n", elseExpr))
	}

	return sb.String()
}

// toFlowchartEvent generates the event part of a flowchart
func (m *Mermaid) toFlowchartEvent(event common.YAMLDocument) string {
	// This is a placeholder as per the Python implementation
	return ""
}

// toFlowchartException generates the exception part of a flowchart
func (m *Mermaid) toFlowchartException(steps []interface{}, exception common.YAMLDocument) string {
	var sb strings.Builder

	for _, s := range steps {
		step := s.(common.YAMLDocument)
		if agent, ok := step["agent"]; ok {
			agentL := m.fixAgentName(agent.(string))
			exceptionName := exception["name"].(string)
			exceptionAgent := exception["agent"].(string)
			sb.WriteString(fmt.Sprintf("%s -->|exception| %s{%s}\n", agentL, exceptionName, exceptionAgent))
		}
	}

	return sb.String()
}
