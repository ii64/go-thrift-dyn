namespace go base

struct Model {
    1: string abc
    4: i64 sd
    9: double f64
}

struct Request {
    6: Model model
    44: list<Model> models
    31: map<i64, Model> modelById
    56: optional set<Model> modset
}

struct ListOnly {
    91: list<string> names
}

struct MapOnly {
    77: map<i64, string> stringById
    88: map<i64, Model> modelById
}

struct EmptyField {}

struct SimpleListMap {
    555: list<i64> listI64
    2016: map<i64, i64> i64ByI64
}

exception BusinessException {
    2: string message
}

service Example {

    Request echoRequest(
        2: Request request) throws (1: BusinessException e)

    oneway void pushAnalytics(
        3: Request request)

}