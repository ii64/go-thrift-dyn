namespace go base

struct Model {
    1: string abc
    4: i64 sd
    9: double f64
    10: optional list<i64> listI64
    11: optional map<i64, i64> mapI64
    12: optional map<i32, i32> mapI32
}

struct Request {
    6: Model model
    7: optional Model model2
    44: list<Model> models
    31: map<i64, Model> modelById
    56: optional set<Model> modset
    88: map<i64, list<Model>> modelByTime
}

struct ListOnly {
    91: list<string> names
}

struct MapOnly {
    10: map<i64, list<i32>> listById
    77: map<i64, string> stringById
    88: map<i64, Model> modelById
}

struct EmptyField {}

struct SimpleListMap {
    10: list<list<i32>> listListI32
    555: list<i64> listI64
    2016: map<i64, i64> i64ByI64

    2: optional list<string> listString
}

struct Common {
    1: binary bin
    2: string bin2
    3: list<i8> bin3
    4: list<byte> bin4
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