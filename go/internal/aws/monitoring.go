package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	costtypes "github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

// MonitoringManager handles AWS monitoring and metrics
type MonitoringManager struct {
	client *Client
}

// NewMonitoringManager creates a new monitoring manager
func NewMonitoringManager(client *Client) *MonitoringManager {
	return &MonitoringManager{
		client: client,
	}
}

// MetricDataPoint represents a single metric data point
type MetricDataPoint struct {
	Timestamp time.Time
	Value     float64
	Unit      string
}

// InstanceMetrics contains metrics for an EC2 instance
type InstanceMetrics struct {
	InstanceID        string
	CPUUtilization    []MetricDataPoint
	NetworkIn         []MetricDataPoint
	NetworkOut        []MetricDataPoint
	DiskReadOps       []MetricDataPoint
	DiskWriteOps      []MetricDataPoint
	StatusCheck       []MetricDataPoint
}

// GetInstanceMetrics retrieves CloudWatch metrics for an EC2 instance
func (mm *MonitoringManager) GetInstanceMetrics(ctx context.Context, instanceID string, startTime, endTime time.Time) (*InstanceMetrics, error) {
	metrics := &InstanceMetrics{
		InstanceID: instanceID,
	}

	// Define metrics to collect
	metricQueries := []struct {
		metricName string
		unit       types.StandardUnit
		statistic  types.Statistic
		result     *[]MetricDataPoint
	}{
		{"CPUUtilization", types.StandardUnitPercent, types.StatisticAverage, &metrics.CPUUtilization},
		{"NetworkIn", types.StandardUnitBytes, types.StatisticSum, &metrics.NetworkIn},
		{"NetworkOut", types.StandardUnitBytes, types.StatisticSum, &metrics.NetworkOut},
		{"DiskReadOps", types.StandardUnitCount, types.StatisticSum, &metrics.DiskReadOps},
		{"DiskWriteOps", types.StandardUnitCount, types.StatisticSum, &metrics.DiskWriteOps},
		{"StatusCheck", types.StandardUnitCount, types.StatisticMaximum, &metrics.StatusCheck},
	}

	for _, query := range metricQueries {
		dataPoints, err := mm.getMetricStatistics(ctx, instanceID, query.metricName, query.statistic, startTime, endTime)
		if err != nil {
			return nil, fmt.Errorf("failed to get %s metrics: %w", query.metricName, err)
		}

		*query.result = dataPoints
	}

	return metrics, nil
}

// getMetricStatistics retrieves metric statistics from CloudWatch
func (mm *MonitoringManager) getMetricStatistics(ctx context.Context, instanceID, metricName string, statistic types.Statistic, startTime, endTime time.Time) ([]MetricDataPoint, error) {
	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String(metricName),
		Dimensions: []types.Dimension{
			{
				Name:  aws.String("InstanceId"),
				Value: aws.String(instanceID),
			},
		},
		StartTime:  aws.Time(startTime),
		EndTime:    aws.Time(endTime),
		Period:     aws.Int32(300), // 5-minute intervals
		Statistics: []types.Statistic{statistic},
	}

	result, err := mm.client.CloudWatch.GetMetricStatistics(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get metric statistics: %w", err)
	}

	dataPoints := make([]MetricDataPoint, 0, len(result.Datapoints))
	for _, dp := range result.Datapoints {
		var value float64
		switch statistic {
		case types.StatisticAverage:
			if dp.Average != nil {
				value = *dp.Average
			}
		case types.StatisticSum:
			if dp.Sum != nil {
				value = *dp.Sum
			}
		case types.StatisticMaximum:
			if dp.Maximum != nil {
				value = *dp.Maximum
			}
		case types.StatisticMinimum:
			if dp.Minimum != nil {
				value = *dp.Minimum
			}
		}

		dataPoint := MetricDataPoint{
			Timestamp: *dp.Timestamp,
			Value:     value,
			Unit:      string(dp.Unit),
		}

		dataPoints = append(dataPoints, dataPoint)
	}

	return dataPoints, nil
}

// CostData represents cost information
type CostData struct {
	Service     string
	Amount      float64
	Unit        string
	Currency    string
	StartDate   string
	EndDate     string
}

// GetCostData retrieves cost data from Cost Explorer
func (mm *MonitoringManager) GetCostData(ctx context.Context, startDate, endDate string, groupBy []string) ([]CostData, error) {
	// Build group by dimensions
	var groupDimensions []costtypes.GroupDefinition
	for _, group := range groupBy {
		groupDimensions = append(groupDimensions, costtypes.GroupDefinition{
			Type: costtypes.GroupDefinitionTypeDimension,
			Key:  aws.String(group),
		})
	}

	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &costtypes.DateInterval{
			Start: aws.String(startDate),
			End:   aws.String(endDate),
		},
		Granularity: costtypes.GranularityDaily,
		Metrics:     []string{"BlendedCost", "UsageQuantity"},
		GroupBy:     groupDimensions,
	}

	result, err := mm.client.CostExplorer.GetCostAndUsage(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get cost and usage data: %w", err)
	}

	var costData []CostData
	for _, resultEntry := range result.ResultsByTime {
		for _, group := range resultEntry.Groups {
			var serviceName string
			if len(group.Keys) > 0 {
				serviceName = group.Keys[0]
			}

			if group.Metrics != nil {
				if blendedCost, exists := group.Metrics["BlendedCost"]; exists && blendedCost.Amount != nil {
					amount := 0.0
					fmt.Sscanf(*blendedCost.Amount, "%f", &amount)

					costEntry := CostData{
						Service:   serviceName,
						Amount:    amount,
						Unit:      "USD", // Cost Explorer typically returns USD
						Currency:  "USD",
						StartDate: *resultEntry.TimePeriod.Start,
						EndDate:   *resultEntry.TimePeriod.End,
					}

					if blendedCost.Unit != nil {
						costEntry.Currency = *blendedCost.Unit
					}

					costData = append(costData, costEntry)
				}
			}
		}
	}

	return costData, nil
}

// AlarmInfo represents CloudWatch alarm information
type AlarmInfo struct {
	AlarmName        string
	AlarmDescription string
	MetricName       string
	Namespace        string
	Statistic        string
	Threshold        float64
	ComparisonOp     string
	State            string
	StateReason      string
	ActionsEnabled   bool
}

// CreateAlarm creates a CloudWatch alarm
func (mm *MonitoringManager) CreateAlarm(ctx context.Context, alarmName, description, metricName, namespace string, threshold float64, comparisonOp string, instanceID string) error {
	input := &cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String(alarmName),
		AlarmDescription:   aws.String(description),
		MetricName:         aws.String(metricName),
		Namespace:          aws.String(namespace),
		Statistic:          types.StatisticAverage,
		Threshold:          aws.Float64(threshold),
		ComparisonOperator: types.ComparisonOperator(comparisonOp),
		EvaluationPeriods:  aws.Int32(2),
		Period:             aws.Int32(300), // 5 minutes
		ActionsEnabled:     aws.Bool(true),
		Dimensions: []types.Dimension{
			{
				Name:  aws.String("InstanceId"),
				Value: aws.String(instanceID),
			},
		},
	}

	_, err := mm.client.CloudWatch.PutMetricAlarm(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create alarm: %w", err)
	}

	return nil
}

// ListAlarms lists CloudWatch alarms
func (mm *MonitoringManager) ListAlarms(ctx context.Context, alarmNamePrefix string) ([]AlarmInfo, error) {
	input := &cloudwatch.DescribeAlarmsInput{}
	
	if alarmNamePrefix != "" {
		input.AlarmNamePrefix = aws.String(alarmNamePrefix)
	}

	result, err := mm.client.CloudWatch.DescribeAlarms(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe alarms: %w", err)
	}

	alarms := make([]AlarmInfo, 0, len(result.MetricAlarms))
	for _, alarm := range result.MetricAlarms {
		alarmInfo := AlarmInfo{
			AlarmName:        *alarm.AlarmName,
			MetricName:       *alarm.MetricName,
			Namespace:        *alarm.Namespace,
			Statistic:        string(alarm.Statistic),
			Threshold:        *alarm.Threshold,
			ComparisonOp:     string(alarm.ComparisonOperator),
			State:            string(alarm.StateValue),
			ActionsEnabled:   *alarm.ActionsEnabled,
		}

		if alarm.AlarmDescription != nil {
			alarmInfo.AlarmDescription = *alarm.AlarmDescription
		}

		if alarm.StateReason != nil {
			alarmInfo.StateReason = *alarm.StateReason
		}

		alarms = append(alarms, alarmInfo)
	}

	return alarms, nil
}

// DeleteAlarm deletes a CloudWatch alarm
func (mm *MonitoringManager) DeleteAlarm(ctx context.Context, alarmName string) error {
	input := &cloudwatch.DeleteAlarmsInput{
		AlarmNames: []string{alarmName},
	}

	_, err := mm.client.CloudWatch.DeleteAlarms(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete alarm: %w", err)
	}

	return nil
}