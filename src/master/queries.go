package master

import (
	"github.com/eaigner/hood"
)

func (rdd *Rdd) GetSegments(tx *hood.Hood) []*Segment {
	var results []Segment
	err := tx.Where("rdd_id", "=", rdd.Id).Find(&results)
	if err != nil {
		panic(err)
	}

	// Should return pointers to the result objects so that
	// they can be mutated
	pointerResults := make([]*Segment, len(results))
	for i := range results {
		pointerResults[i] = &results[i]
	}

	return pointerResults
}

func (workflow *Workflow) GetProtojobs(tx *hood.Hood) []*Protojob {
	var results []Protojob
	err := tx.Where("workflow_id", "=", workflow.Id).Find(&results)
	if err != nil {
		panic(err)
	}

	// Should return pointers to the result objects so that
	// they can be mutated
	pointerResults := make([]*Protojob, len(results))
	for i := range results {
		pointerResults[i] = &results[i]
	}

	return pointerResults
}

// Select all workflow edges whose dest_protojob is in the given
// workflow (this also imlies that the source_protojob is in the workflow)
func (workflow *Workflow) GetWorkflowEdges(tx *hood.Hood) []*WorkflowEdge {
	var results []WorkflowEdge
	err := tx.FindSql(&results,
		`select *
    from workflow_edge
    inner join protojob dest_job
    on workflow_edge.dest_job_id = dest_job.id
    where dest_job.workflow_id = $1`, workflow.Id)
	if err != nil {
		panic(err)
	}

	// Should return pointers to the result objects so that
	// they can be mutated
	pointerResults := make([]*WorkflowEdge, len(results))
	for i := range results {
		pointerResults[i] = &results[i]
	}

	return pointerResults
}

// Get all Rdd edges whose dest Rdd is in the given workflow batch. Note that
// the source and dest Rdds may be produced in different batches.
func (workflowBatch *WorkflowBatch) GetRddEdges(tx *hood.Hood) []*RddEdge {
	var results []RddEdge
	err := tx.FindSql(&results,
		`select *
    from rdd_edge
    inner join rdd dest_rdd
    on rdd_edge.dest_rdd_id = dest_rdd.id
    where dest_rdd.workflow_batch_id = $1`, workflowBatch.Id)
	if err != nil {
		panic(err)
	}

	// Should return pointers to the result objects so that
	// they can be mutated
	pointerResults := make([]*RddEdge, len(results))
	for i := range results {
		pointerResults[i] = &results[i]
	}

	return pointerResults
}

// Same as GetRddEdges, except the edges cannot have any delay
func (workflowBatch *WorkflowBatch) GetNonDelayRddEdges(tx *hood.Hood) []*RddEdge {
	var results []RddEdge
	err := tx.FindSql(&results,
		`select *
    from rdd_edge
    inner join rdd dest_rdd
    on rdd_edge.dest_rdd_id = dest_rdd.id
    inner join workflow_edge
    on workflow_edge.id = rdd_edge.workflow_edge_id
    where dest_rdd.workflow_batch_id = $1
    and workflow_edge.delay=0`, workflowBatch.Id)
	if err != nil {
		panic(err)
	}

	// Should return pointers to the result objects so that
	// they can be mutated
	pointerResults := make([]*RddEdge, len(results))
	for i := range results {
		pointerResults[i] = &results[i]
	}

	return pointerResults
}

func (workflow *Workflow) GetWorkflowBatches(tx *hood.Hood) []*WorkflowBatch {
	var results []WorkflowBatch
	err := tx.Where("workflow_id", "=", workflow.Id).Find(&results)
	if err != nil {
		panic(err)
	}

	// Should return pointers to the result objects so that
	// they can be mutated
	pointerResults := make([]*WorkflowBatch, len(results))
	for i := range results {
		pointerResults[i] = &results[i]
	}

	return pointerResults
}

func (workflowBatch *WorkflowBatch) GetRdds(tx *hood.Hood) []*Rdd {
	var results []Rdd
	err := tx.Where("workflow_batch_id", "=", workflowBatch.Id).Find(&results)
	if err != nil {
		panic(err)
	}

	// Should return pointers to the result objects so that
	// they can be mutated
	pointerResults := make([]*Rdd, len(results))
	for i := range results {
		pointerResults[i] = &results[i]
	}

	return pointerResults
}

func GetWorkflows(tx *hood.Hood) []*Workflow {
	var results []Workflow
	err := tx.Find(&results)
	if err != nil {
		panic(err)
	}

	// Should return pointers to the result objects so that
	// they can be mutated
	pointerResults := make([]*Workflow, len(results))
	for i := range results {
		pointerResults[i] = &results[i]
	}

	return pointerResults
}

func GetWorkers(tx *hood.Hood) []*Worker {
	var results []Worker
	err := tx.Find(&results)
	if err != nil {
		panic(err)
	}

	// Should return pointers to the result objects so that
	// they can be mutated
	pointerResults := make([]*Worker, len(results))
	for i := range results {
		pointerResults[i] = &results[i]
	}

	return pointerResults
}

func GetWorkersAtAddress(tx *hood.Hood, address string) []*Worker {
	var results []Worker
	err := tx.Where("url", "=", address).Find(&results)
	if err != nil {
		panic(err)
	}

	// Should return pointers to the result objects so that
	// they can be mutated
	pointerResults := make([]*Worker, len(results))
	for i := range results {
		pointerResults[i] = &results[i]
	}

	return pointerResults
}

func GetWorker(tx *hood.Hood, id int64) *Worker {
	var results []Worker
	err := tx.Where("id", "=", id).Find(&results)
	if err != nil {
		panic(err)
	}

	if len(results) == 0 {
		return nil
	} else {
		return &results[0]
	}
}

func GetRddByStartTime(tx *hood.Hood, protojobId int64, startTime int) *Rdd {
	var results []Rdd
	err := tx.Join(hood.InnerJoin, &WorkflowBatch{}, "workflow_batch.id", "rdd.workflow_batch_id").Where("protojob_id", "=", protojobId).And("start_time", "=", startTime).Find(&results)
	if err != nil {
		panic(err)
	}

	if len(results) == 0 {
		return nil
	} else {
		return &results[0]
	}
}

func GetLastWorkflowBatch(tx *hood.Hood, workflow *Workflow) *WorkflowBatch {
	var results []WorkflowBatch
	err := tx.Where("workflow_id", "=", workflow.Id).OrderBy("start_time").Desc().Limit(1).Find(&results)
	if err != nil {
		panic(err)
	}

	if len(results) == 0 {
		return nil
	} else {
		return &results[0]
	}
}
