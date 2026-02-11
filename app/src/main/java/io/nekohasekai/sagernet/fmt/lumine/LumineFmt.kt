package io.nekohasekai.sagernet.fmt.lumine

import io.nekohasekai.sagernet.ktx.toStringPretty
import org.json.JSONObject

fun LumineBean.buildLumineConfig(): String {
    return JSONObject().apply {
        put("socks5_address", socks5Address)
        put("http_address", httpAddress)
        put("dns_addr", dnsAddr)
        put("dns_cache_ttl", dnsCacheTTL)
        put("default_policy", JSONObject().apply {
            put("mode", "proxy")
            put("connect_timeout", "10s")
            put("fragment_size", fragmentSize)
            put("fragment_sleep", fragmentSleep)
            put("fragment_type", fragmentType)
            put("desync_split", desyncSplit)
            put("desync_split_position", desyncSplitPosition)
        })
        put("domain_policies", JSONObject())
        put("ip_policies", JSONObject())
    }.toStringPretty()
}
