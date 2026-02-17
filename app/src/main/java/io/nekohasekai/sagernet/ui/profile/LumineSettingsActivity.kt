package io.nekohasekai.sagernet.ui.profile

import android.os.Bundle
import androidx.preference.EditTextPreference
import androidx.preference.ListPreference
import androidx.preference.PreferenceFragmentCompat
import io.nekohasekai.sagernet.Key
import io.nekohasekai.sagernet.R
import io.nekohasekai.sagernet.database.DataStore
import io.nekohasekai.sagernet.database.preference.EditTextPreferenceModifiers
import io.nekohasekai.sagernet.fmt.lumine.LumineBean

class LumineSettingsActivity : ProfileSettingsActivity<LumineBean>() {

    override fun createEntity() = LumineBean()

    override fun LumineBean.init() {
        DataStore.profileName = name
        DataStore.serverAddress = socks5Address.substringBefore(":")
        DataStore.serverPort = socks5Address.substringAfter(":", "1080").toIntOrNull() ?: 1080
        
        DataStore.profileCacheStore.putString("lumine_socks5_address", socks5Address)
        DataStore.profileCacheStore.putString("lumine_http_address", httpAddress)
        DataStore.profileCacheStore.putString("lumine_dns_addr", dnsAddr)
        DataStore.profileCacheStore.putInt("lumine_dns_cache_ttl", dnsCacheTTL ?: 3600)
        DataStore.profileCacheStore.putInt("lumine_fragment_size", fragmentSize ?: 1024)
        DataStore.profileCacheStore.putString("lumine_fragment_sleep", fragmentSleep)
        DataStore.profileCacheStore.putString("lumine_fragment_type", fragmentType)
        DataStore.profileCacheStore.putInt("lumine_desync_split", desyncSplit ?: 2)
        DataStore.profileCacheStore.putInt("lumine_desync_split_position", desyncSplitPosition ?: 3)
    }

    override fun LumineBean.serialize() {
        name = DataStore.profileName
        
        socks5Address = DataStore.profileCacheStore.getString("lumine_socks5_address") ?: "127.0.0.1:1080"
        httpAddress = DataStore.profileCacheStore.getString("lumine_http_address") ?: "none"
        dnsAddr = DataStore.profileCacheStore.getString("lumine_dns_addr") ?: "https://1.1.1.1/dns-query"
        dnsCacheTTL = DataStore.profileCacheStore.getInt("lumine_dns_cache_ttl")
        fragmentSize = DataStore.profileCacheStore.getInt("lumine_fragment_size")
        fragmentSleep = DataStore.profileCacheStore.getString("lumine_fragment_sleep") ?: "10ms"
        fragmentType = DataStore.profileCacheStore.getString("lumine_fragment_type") ?: "tls"
        desyncSplit = DataStore.profileCacheStore.getInt("lumine_desync_split")
        desyncSplitPosition = DataStore.profileCacheStore.getInt("lumine_desync_split_position")
        
        // 更新 serverAddress 和 serverPort 用于显示
        serverAddress = socks5Address.substringBefore(":")
        serverPort = socks5Address.substringAfter(":", "1080").toIntOrNull() ?: 1080
    }

    override fun PreferenceFragmentCompat.createPreferences(
        savedInstanceState: Bundle?,
        rootKey: String?,
    ) {
        addPreferencesFromResource(R.xml.lumine_preferences)
        
        findPreference<EditTextPreference>("lumine_socks5_address")!!.apply {
            setOnBindEditTextListener(EditTextPreferenceModifiers.Hosts)
        }
        
        findPreference<EditTextPreference>("lumine_dns_cache_ttl")!!.apply {
            setOnBindEditTextListener(EditTextPreferenceModifiers.Number)
        }
        
        findPreference<EditTextPreference>("lumine_fragment_size")!!.apply {
            setOnBindEditTextListener(EditTextPreferenceModifiers.Number)
        }
        
        findPreference<EditTextPreference>("lumine_desync_split")!!.apply {
            setOnBindEditTextListener(EditTextPreferenceModifiers.Number)
        }
        
        findPreference<EditTextPreference>("lumine_desync_split_position")!!.apply {
            setOnBindEditTextListener(EditTextPreferenceModifiers.Number)
        }
    }
}
