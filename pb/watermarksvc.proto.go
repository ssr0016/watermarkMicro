package pb

syntax = "proto3";

package pb;

service WatermarkSevice {
	rpc Get(GetRequest) returns (GetReply) {}

	rpc Watermark (WatermarkRequest) returns (WatermarkReply) {}

	rpc status (StatusRequest) returns (StatusReply) {}

	rpc AddDocument (AddDocumentRequest) returns (AddDocumentReply) {}

	rpc ServiceStatus (ServiceStatusRequest) returns (ServiceStatusReply) {}
}

message Document {
	string content = 1;
	string title = 2;
	string author = 3;
	string topic = 4;
	string watermark = 5;
}

message GetRequest {
	repeated Filters {
		string key = 1;
		string value = 2;
	}
	repeated Filters = 1;
}

message GetReply {
	repeated Document document = 1
	string Err = 2;
}

message StatusRequest {
	string ticketID = 1;
}

message StatusReply {
	enum Status {
		PENDING = 0;
		STARTED = 1;
		IN_PROGRESS = 2;
		FINISHED = 3;
		FAILED = 4;
	}
	Status status = 1;
	string Err = 2;
}

message WatermarkRequest {
	string ticketID = 1;
	string mark = 2;
}

message WatermarkReply{
	int64 code = 1;
	string err = 2;
}

message AddDocumentRequest {
	Document document = 1;
}

message AddDocumentReply{
	string ticketID = 1;
	string err = 2;
}

message ServiceStatusRequest{}

message ServiceStatusReply{
	int64 code = 1;
	string err = 2;
}


