package main

import (
   "flag"
   "fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
   "time"
)

type Metric struct {
   MetricName string
   Namespace string
   Statistics []string
   Dimensions []*cloudwatch.Dimension
   Unit string
}

type Request struct {
   Period int64
   StartTime time.Time
   EndTime time.Time
   MetricList []Metric
}

func parseDate(instr string, default_offset int64) (time.Time, error) {
   retDate := time.Now()
   if instr != "" {
      // Try to see if it is a duration first
      offset, err := time.ParseDuration(instr)
      if err != nil {
         // Now try RFC 3339 format
         
      }
    
   } else {
      
   }
}

func parseArgs() {
   r := &Request{}
   flag.Int64Var(&r.Period, "period", 60, "Number of seconds in each period")
   startTime := flag.String("start_time", "", "Start time in RFC 3339 format")
   endTime := flag.String("end_time", "", "End time in RFC 3339 format")

   metricFile := flag.String("metric_file", "metrics.json", "The filename of the desired metrics to be queried")

   flag.Parse()

   if startTime == "" {
      r.StartTime = time.Now().Add(time.Duration(-60*time.Second))
   } else {
   }


}

func main() {
	svc := cloudwatch.New(session.New())

	params := &cloudwatch.GetMetricStatisticsInput{
		EndTime:    aws.Time(time.Now()),     // Required
		MetricName: aws.String("MetricName"), // Required
		Namespace:  aws.String("Namespace"),  // Required
		Period:     aws.Int64(1),             // Required
		StartTime:  aws.Time(time.Now()),     // Required
		Statistics: []*string{ // Required
			aws.String("Statistic"), // Required
			// More values...
		},
		Dimensions: []*cloudwatch.Dimension{
			{ // Required
				Name:  aws.String("DimensionName"),  // Required
				Value: aws.String("DimensionValue"), // Required
			},
			// More values...
		},
		Unit: aws.String("StandardUnit"),
	}
	resp, err := svc.GetMetricStatistics(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}
